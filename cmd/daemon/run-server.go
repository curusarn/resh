package main

import (
	"net/http"
	"strconv"

	"github.com/curusarn/resh/pkg/cfg"
	"github.com/curusarn/resh/pkg/histfile"
	"github.com/curusarn/resh/pkg/records"
	"github.com/curusarn/resh/pkg/sesshist"
	"github.com/curusarn/resh/pkg/sesswatch"
)

func runServer(config cfg.Config, outputPath string) {
	var recordSubscribers []chan records.Record
	var sessionInitSubscribers []chan records.Record
	var sessionDropSubscribers []chan string

	// sessshist
	sesshistSessionsToInit := make(chan records.Record)
	sessionInitSubscribers = append(sessionInitSubscribers, sesshistSessionsToInit)
	sesshistSessionsToDrop := make(chan string)
	sessionDropSubscribers = append(sessionDropSubscribers, sesshistSessionsToDrop)
	sesshistRecords := make(chan records.Record)
	recordSubscribers = append(recordSubscribers, sesshistRecords)
	sesshistDispatch := sesshist.NewDispatch(sesshistSessionsToInit, sesshistSessionsToDrop, sesshistRecords)

	// histfile
	histfileRecords := make(chan records.Record)
	recordSubscribers = append(recordSubscribers, histfileRecords)
	histfileSessionsToDrop := make(chan string)
	sessionDropSubscribers = append(sessionDropSubscribers, histfileSessionsToDrop)
	histfile.Go(histfileRecords, outputPath, histfileSessionsToDrop)

	// sesswatch
	sesswatchSessionsToWatch := make(chan records.Record)
	sessionInitSubscribers = append(sessionInitSubscribers, sesswatchSessionsToWatch)
	sesswatch.Go(sesswatchSessionsToWatch, sessionDropSubscribers, config.SesswatchPeriodSeconds)

	// handlers
	http.HandleFunc("/status", statusHandler)
	http.Handle("/record", &recordHandler{subscribers: recordSubscribers})
	http.Handle("/session_init", &sessionInitHandler{subscribers: sessionInitSubscribers})
	http.Handle("/recall", &recallHandler{sesshistDispatch: sesshistDispatch})
	http.ListenAndServe(":"+strconv.Itoa(config.Port), nil)
}
