package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/curusarn/resh/internal/recordint"
	"go.uber.org/zap"
)

func NewRecordHandler(sugar *zap.SugaredLogger, subscribers []chan recordint.Collect) recordHandler {
	return recordHandler{
		sugar:       sugar.With(zap.String("endpoint", "/record")),
		subscribers: subscribers,
	}
}

type recordHandler struct {
	sugar       *zap.SugaredLogger
	subscribers []chan recordint.Collect
}

func (h *recordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sugar := h.sugar.With(zap.String("endpoint", "/record"))
	sugar.Debugw("Handling request, sending response, reading body ...")
	w.Write([]byte("OK\n"))
	jsn, err := ioutil.ReadAll(r.Body)
	// run rest of the handler as goroutine to prevent any hangups
	go func() {
		if err != nil {
			sugar.Errorw("Error reading body", "error", err)
			return
		}

		sugar.Debugw("Unmarshaling record ...")
		rec := recordint.Collect{}
		err = json.Unmarshal(jsn, &rec)
		if err != nil {
			sugar.Errorw("Error during unmarshaling",
				"error", err,
				"payload", jsn,
			)
			return
		}
		part := "2"
		if rec.Rec.PartOne {
			part = "1"
		}
		sugar := sugar.With(
			"cmdLine", rec.Rec.CmdLine,
			"part", part,
		)
		sugar.Debugw("Got record, sending to subscribers ...")
		for _, sub := range h.subscribers {
			sub <- rec
		}
		sugar.Debugw("Record sent to subscribers")
	}()
}
