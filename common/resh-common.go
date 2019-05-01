package common

import (
    "bytes"
    "encoding/json"
    "net/http"
    "log"
)

type Record struct {
    CmdLine string     `json: cmdLine`
    Pwd string         `json: pwd`
    GitWorkTree string `json: gitWorkTree`
    Shell string       `json: shell`
    ExitCode int       `json: exitCode`
    Logs []string      `json: logs`
}

func (r Record) Send() {
    recJson, err := json.Marshal(r)
    if err != nil {
        log.Fatal("1", err)
    }

    req, err := http.NewRequest("POST", "http://localhost:8888",
                                bytes.NewBuffer(recJson))
    if err != nil {
        log.Fatal("2", err)
    }
	req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    _, err = client.Do(req)
	if err != nil {
		log.Fatal("3", err)
	}
}
