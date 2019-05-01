package main

import (
    "encoding/json"
    "log"
    "io/ioutil"
    "net/http"
    common "github.com/curusarn/resh/common"
)

func main() {
    server()
}

func recordHandler(w http.ResponseWriter, r *http.Request) {
    record := common.Record{}

    jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading the body", err)
        return
	}

    err = json.Unmarshal(jsn, &record)
	if err != nil {
		log.Println("Decoding error: ", err)
        return
	}

    log.Printf("Received: %v\n", record)
    // fmt.Println("cmd:", r.CmdLine)
    // fmt.Println("pwd:", r.Pwd)
    // fmt.Println("git:", r.GitWorkTree)
    // fmt.Println("exit_code:", r.ExitCode)
    w.Write([]byte("OK\n"))
}

func server() {
	http.HandleFunc("/", recordHandler)
	http.ListenAndServe(":8888", nil)
}
