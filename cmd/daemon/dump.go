package main

import (
	"encoding/json"
	"github.com/curusarn/resh/internal/histcli"
	"io"
	"net/http"

	"github.com/curusarn/resh/internal/msg"
	"go.uber.org/zap"
)

type dumpHandler struct {
	sugar   *zap.SugaredLogger
	history *histcli.Histcli
}

func (h *dumpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sugar := h.sugar.With(zap.String("endpoint", "/dump"))
	sugar.Debugw("Handling request, reading body ...")
	jsn, err := io.ReadAll(r.Body)
	if err != nil {
		sugar.Errorw("Error reading body", "error", err)
		return
	}

	sugar.Debugw("Unmarshalling record ...")
	mess := msg.CliMsg{}
	err = json.Unmarshal(jsn, &mess)
	if err != nil {
		sugar.Errorw("Error during unmarshalling",
			"error", err,
			"payload", jsn,
		)
		return
	}
	sugar.Debugw("Getting records to send ...")

	resp := msg.CliResponse{Records: h.history.Dump()}
	jsn, err = json.Marshal(&resp)
	if err != nil {
		sugar.Errorw("Error when marshaling", "error", err)
		return
	}
	w.Write(jsn)
	sugar.Infow("Request handled")
}
