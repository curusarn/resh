package main

import (
	"encoding/json"
	"net/http"

	"github.com/curusarn/resh/internal/msg"
	"go.uber.org/zap"
)

type statusHandler struct {
	sugar *zap.SugaredLogger
}

func (h *statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sugar := h.sugar.With(zap.String("endpoint", "/status"))
	sugar.Debugw("Handling request ...")
	resp := msg.StatusResponse{
		Status:  true,
		Version: version,
		Commit:  commit,
	}
	jsn, err := json.Marshal(&resp)
	if err != nil {
		sugar.Errorw("Error when marshaling",
			"error", err,
			"response", resp,
		)
		return
	}
	w.Write(jsn)
	sugar.Infow("Request handled")
}
