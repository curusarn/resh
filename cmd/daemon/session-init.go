package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/curusarn/resh/internal/recordint"
	"go.uber.org/zap"
)

type sessionInitHandler struct {
	sugar       *zap.SugaredLogger
	subscribers []chan recordint.SessionInit
}

func (h *sessionInitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sugar := h.sugar.With(zap.String("endpoint", "/session_init"))
	sugar.Debugw("Handling request, sending response, reading body ...")
	w.Write([]byte("OK\n"))
	// TODO: should we somehow check for errors here?
	jsn, err := io.ReadAll(r.Body)
	// run rest of the handler as goroutine to prevent any hangups
	go func() {
		if err != nil {
			sugar.Errorw("Error reading body", "error", err)
			return
		}

		sugar.Debugw("Unmarshalling record ...")
		rec := recordint.SessionInit{}
		err = json.Unmarshal(jsn, &rec)
		if err != nil {
			sugar.Errorw("Error during unmarshalling",
				"error", err,
				"payload", jsn,
			)
			return
		}
		sugar := sugar.With(
			"sessionID", rec.SessionID,
			"sessionPID", rec.SessionPID,
		)
		sugar.Infow("Got session, sending to subscribers ...")
		for _, sub := range h.subscribers {
			sub <- rec
		}
		sugar.Debugw("Session sent to subscribers")
	}()
}
