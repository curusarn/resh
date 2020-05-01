package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/curusarn/resh/pkg/histfile"
	"github.com/curusarn/resh/pkg/msg"
)

type dumpHandler struct {
	histfileBox *histfile.Histfile
}

func (h *dumpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if Debug {
		log.Println("/dump START")
		log.Println("/dump reading body ...")
	}
	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading the body", err)
		return
	}

	mess := msg.CliMsg{}
	if Debug {
		log.Println("/dump unmarshaling record ...")
	}
	err = json.Unmarshal(jsn, &mess)
	if err != nil {
		log.Println("Decoding error:", err)
		log.Println("Payload:", jsn)
		return
	}
	if Debug {
		log.Println("/dump dumping ...")
	}
	fullRecords := h.histfileBox.DumpCliRecords()
	if err != nil {
		log.Println("Dump error:", err)
	}

	resp := msg.CliResponse{CliRecords: fullRecords.List}
	jsn, err = json.Marshal(&resp)
	if err != nil {
		log.Println("Encoding error:", err)
		return
	}
	w.Write(jsn)
	log.Println("/dump END")
}
