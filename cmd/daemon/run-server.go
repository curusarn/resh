package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/curusarn/resh/internal/cfg"
	"github.com/curusarn/resh/internal/histfile"
	"github.com/curusarn/resh/internal/recordint"
	"github.com/curusarn/resh/internal/sesswatch"
	"github.com/curusarn/resh/internal/signalhandler"
	"go.uber.org/zap"
)

// TODO: turn server and handlers into package

type Server struct {
	sugar  *zap.SugaredLogger
	config cfg.Config

	reshHistoryPath string
	bashHistoryPath string
	zshHistoryPath  string
}

func (s *Server) Run() {
	var recordSubscribers []chan recordint.Collect
	var sessionInitSubscribers []chan recordint.SessionInit
	var sessionDropSubscribers []chan string
	var signalSubscribers []chan os.Signal

	shutdown := make(chan string)

	// histfile
	histfileRecords := make(chan recordint.Collect)
	recordSubscribers = append(recordSubscribers, histfileRecords)
	histfileSessionsToDrop := make(chan string)
	sessionDropSubscribers = append(sessionDropSubscribers, histfileSessionsToDrop)
	histfileSignals := make(chan os.Signal)
	signalSubscribers = append(signalSubscribers, histfileSignals)
	maxHistSize := 10000  // lines
	minHistSizeKB := 2000 // roughly lines
	histfileBox := histfile.New(s.sugar, histfileRecords, histfileSessionsToDrop,
		s.reshHistoryPath, s.bashHistoryPath, s.zshHistoryPath,
		maxHistSize, minHistSizeKB,
		histfileSignals, shutdown)

	// sesswatch
	sesswatchRecords := make(chan recordint.Collect)
	recordSubscribers = append(recordSubscribers, sesswatchRecords)
	sesswatchSessionsToWatch := make(chan recordint.SessionInit)
	sessionInitSubscribers = append(sessionInitSubscribers, sesswatchSessionsToWatch)
	sesswatch.Go(
		s.sugar,
		sesswatchSessionsToWatch,
		sesswatchRecords,
		sessionDropSubscribers,
		s.config.SesswatchPeriodSeconds,
	)

	// handlers
	mux := http.NewServeMux()
	mux.Handle("/status", &statusHandler{sugar: s.sugar})
	mux.Handle("/record", &recordHandler{sugar: s.sugar, subscribers: recordSubscribers})
	mux.Handle("/session_init", &sessionInitHandler{sugar: s.sugar, subscribers: sessionInitSubscribers})
	mux.Handle("/dump", &dumpHandler{sugar: s.sugar, histfileBox: histfileBox})

	server := &http.Server{
		Addr:              "localhost:" + strconv.Itoa(s.config.Port),
		Handler:           mux,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	go server.ListenAndServe()

	// signalhandler - takes over the main goroutine so when signal handler exists the whole program exits
	signalhandler.Run(s.sugar, signalSubscribers, shutdown, server)
}
