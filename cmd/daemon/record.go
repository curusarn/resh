package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/curusarn/resh/pkg/records"
)

type recordHandler struct {
	subscribers []chan records.Record
}

func (h *recordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		for _, sub := range h.subscribers {
			sub <- record
		}
		part := "2"
		if record.PartOne {
			part = "1"
		}
		log.Println("/record - ", record.CmdLine, " - part", part)
	}()

	// fmt.Println("cmd:", r.CmdLine)
	// fmt.Println("pwd:", r.Pwd)
	// fmt.Println("git:", r.GitWorkTree)
	// fmt.Println("exit_code:", r.ExitCode)
}
