package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/curusarn/resh/internal/histfile"
	"github.com/curusarn/resh/internal/msg"
	"go.uber.org/zap"
)

type dumpHandler struct {
	sugar       *zap.SugaredLogger
	histfileBox *histfile.Histfile
}

func (h *dumpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sugar := h.sugar.With(zap.String("endpoint", "/dump"))
	sugar.Debugw("Handling request, reading body ...")
	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sugar.Errorw("Error reading body", "error", err)
		return
	}

	sugar.Debugw("Unmarshaling record ...")
	mess := msg.CliMsg{}
	err = json.Unmarshal(jsn, &mess)
	if err != nil {
		sugar.Errorw("Error during unmarshaling",
			"error", err,
			"payload", jsn,
		)
		return
	}
	sugar.Debugw("Getting records to send ...")
	fullRecords := h.histfileBox.DumpCliRecords()
	if err != nil {
		sugar.Errorw("Error when getting records", "error", err)
	}

	resp := msg.CliResponse{Records: fullRecords.List}
	jsn, err = json.Marshal(&resp)
	if err != nil {
		sugar.Errorw("Error when marshaling", "error", err)
		return
	}
	w.Write(jsn)
	sugar.Infow("Request handled")
}
