package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/curusarn/resh/pkg/msg"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/status START")
	resp := msg.StatusResponse{
		Status:  true,
		Version: version,
		Commit:  commit,
	}
	jsn, err := json.Marshal(&resp)
	if err != nil {
		log.Println("Encoding error:", err)
		log.Println("Response:", resp)
		return
	}
	w.Write(jsn)
	log.Println("/status END")
}

func isDaemonRunning(port int) (bool, error) {
	url := "http://localhost:" + strconv.Itoa(port) + "/status"
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error while checking daemon status - "+
			"it's probably not running!", err)
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}
