package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/curusarn/resh/pkg/collect"
	"github.com/curusarn/resh/pkg/msg"
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
	found := true
	cmd, err := h.sesshistDispatch.Recall(rec.SessionID, rec.RecallHistno, rec.RecallPrefix)
	if err != nil {
		log.Println("/recall - sess id:", rec.SessionID, " - histno:", rec.RecallHistno, " -> ERROR")
		log.Println("Recall error:", err)
		found = false
		cmd = ""
	}
	resp := collect.SingleResponse{CmdLine: cmd, Found: found}
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
	log.Println("/recall END - sess id:", rec.SessionID, " - histno:", rec.RecallHistno, " -> ", cmd, " (found:", found, ")")
}

type inspectHandler struct {
	sesshistDispatch *sesshist.Dispatch
}

func (h *inspectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("/inspect START")
	log.Println("/inspect reading body ...")
	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading the body", err)
		return
	}

	mess := msg.InspectMsg{}
	log.Println("/inspect unmarshaling record ...")
	err = json.Unmarshal(jsn, &mess)
	if err != nil {
		log.Println("Decoding error:", err)
		log.Println("Payload:", jsn)
		return
	}
	log.Println("/inspect recalling ...")
	cmds, err := h.sesshistDispatch.Inspect(mess.SessionID, int(mess.Count))
	if err != nil {
		log.Println("/inspect - sess id:", mess.SessionID, " - count:", mess.Count, " -> ERROR")
		log.Println("Inspect error:", err)
		return
	}
	resp := msg.MultiResponse{CmdLines: cmds}
	log.Println("/inspect marshaling response ...")
	jsn, err = json.Marshal(&resp)
	if err != nil {
		log.Println("Encoding error:", err)
		log.Println("Response:", resp)
		return
	}
	// log.Println(string(jsn))
	log.Println("/inspect writing response ...")
	w.Write(jsn)
	log.Println("/inspect END - sess id:", mess.SessionID, " - count:", mess.Count)
}
