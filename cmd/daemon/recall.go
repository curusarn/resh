package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/curusarn/resh/pkg/collect"
	"github.com/curusarn/resh/pkg/records"
	"github.com/curusarn/resh/pkg/sesshist"
)

type recallHandler struct {
	sesshistDispatch *sesshist.Dispatch
}

func (h *recallHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading the body", err)
		return
	}

	rec := records.Record{}
	err = json.Unmarshal(jsn, &rec)
	if err != nil {
		log.Println("Decoding error:", err)
		log.Println("Payload:", jsn)
		return
	}
	cmd, err := h.sesshistDispatch.Recall(rec.SessionID, rec.RecallHistno)
	if err != nil {
		log.Println("/recall - sess id:", rec.SessionID, " - histno:", rec.RecallHistno, " -> ERROR")
		log.Println("Recall error:", err)
		return
	}
	resp := collect.SingleResponse{cmd}
	jsn, err = json.Marshal(&resp)
	if err != nil {
		log.Println("Encoding error:", err)
		log.Println("Response:", resp)
		return
	}
	log.Println(string(jsn))
	w.Write(jsn)
	log.Println("/recall - sess id:", rec.SessionID, " - histno:", rec.RecallHistno, " -> ", cmd)
}
