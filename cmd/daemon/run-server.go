package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/curusarn/resh/pkg/cfg"
	"github.com/curusarn/resh/pkg/histfile"
	"github.com/curusarn/resh/pkg/records"
	"github.com/curusarn/resh/pkg/sesshist"
	"github.com/curusarn/resh/pkg/sesswatch"
	"github.com/curusarn/resh/pkg/signalhandler"
)

func runServer(config cfg.Config, reshHistoryPath, bashHistoryPath, zshHistoryPath string) {
	var recordSubscribers []chan records.Record
	var sessionInitSubscribers []chan records.Record
	var sessionDropSubscribers []chan string
	var signalSubscribers []chan os.Signal

	shutdown := make(chan string)

	// sessshist
	sesshistSessionsToInit := make(chan records.Record)
	sessionInitSubscribers = append(sessionInitSubscribers, sesshistSessionsToInit)
	sesshistSessionsToDrop := make(chan string)
	sessionDropSubscribers = append(sessionDropSubscribers, sesshistSessionsToDrop)
	sesshistRecords := make(chan records.Record)
	recordSubscribers = append(recordSubscribers, sesshistRecords)

	// histfile
	histfileRecords := make(chan records.Record)
	recordSubscribers = append(recordSubscribers, histfileRecords)
	histfileSessionsToDrop := make(chan string)
	sessionDropSubscribers = append(sessionDropSubscribers, histfileSessionsToDrop)
	histfileSignals := make(chan os.Signal)
	signalSubscribers = append(signalSubscribers, histfileSignals)
	maxHistSize := 10000  // lines
	minHistSizeKB := 2000 // roughly lines
	histfileBox := histfile.New(histfileRecords, histfileSessionsToDrop,
		reshHistoryPath, bashHistoryPath, zshHistoryPath,
		maxHistSize, minHistSizeKB,
		histfileSignals, shutdown)

	// sesshist New
	sesshistDispatch := sesshist.NewDispatch(sesshistSessionsToInit, sesshistSessionsToDrop,
		sesshistRecords, histfileBox,
		config.SesshistInitHistorySize)

	// sesswatch
	sesswatchRecords := make(chan records.Record)
	recordSubscribers = append(recordSubscribers, sesswatchRecords)
	sesswatchSessionsToWatch := make(chan records.Record)
	sessionInitSubscribers = append(sessionInitSubscribers, sesswatchRecords, sesswatchSessionsToWatch)
	sesswatch.Go(sesswatchSessionsToWatch, sesswatchRecords, sessionDropSubscribers, config.SesswatchPeriodSeconds)

	// handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/status", statusHandler)
	mux.Handle("/record", &recordHandler{subscribers: recordSubscribers})
	mux.Handle("/session_init", &sessionInitHandler{subscribers: sessionInitSubscribers})
	mux.Handle("/recall", &recallHandler{sesshistDispatch: sesshistDispatch})
	mux.Handle("/inspect", &inspectHandler{sesshistDispatch: sesshistDispatch})
	mux.Handle("/dump", &dumpHandler{histfileBox: histfileBox})

	server := &http.Server{Addr: "localhost:" + strconv.Itoa(config.Port), Handler: mux}
	go server.ListenAndServe()

	// signalhandler - takes over the main goroutine so when signal handler exists the whole program exits
	signalhandler.Run(signalSubscribers, shutdown, server)
}
