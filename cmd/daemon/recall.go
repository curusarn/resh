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
	log.Println("/recall START")
	log.Println("/recall reading body ...")
	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading the body", err)
		return
	}

	rec := records.SlimRecord{}
	log.Println("/recall unmarshaling record ...")
	err = json.Unmarshal(jsn, &rec)
	if err != nil {
		log.Println("Decoding error:", err)
		log.Println("Payload:", jsn)
		return
	}
	log.Println("/recall recalling ...")
	cmd, err := h.sesshistDispatch.Recall(rec.SessionID, rec.RecallHistno, rec.RecallPrefix)
	if err != nil {
		log.Println("/recall - sess id:", rec.SessionID, " - histno:", rec.RecallHistno, " -> ERROR")
		log.Println("Recall error:", err)
		return
	}
	resp := collect.SingleResponse{CmdLine: cmd}
	log.Println("/recall marshaling response ...")
	jsn, err = json.Marshal(&resp)
	if err != nil {
		log.Println("Encoding error:", err)
		log.Println("Response:", resp)
		return
	}
	log.Println(string(jsn))
	log.Println("/recall writing response ...")
	w.Write(jsn)
	log.Println("/recall END - sess id:", rec.SessionID, " - histno:", rec.RecallHistno, " -> ", cmd)
}
