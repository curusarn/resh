package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/curusarn/resh/pkg/records"
)

type sessionInitHandler struct {
	subscribers []chan records.Record
}

func (h *sessionInitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK\n"))
	jsn, err := ioutil.ReadAll(r.Body)
	// run rest of the handler as goroutine to prevent any hangups
	go func() {
		if err != nil {
			log.Println("Error reading the body", err)
			return
		}

		record := records.Record{}
		err = json.Unmarshal(jsn, &record)
		if err != nil {
			log.Println("Decoding error: ", err)
			log.Println("Payload: ", jsn)
			return
		}
		log.Println("/session_init - id:", record.SessionID, " - pid:", record.SessionPID)
		for _, sub := range h.subscribers {
			sub <- record
		}
	}()
}
