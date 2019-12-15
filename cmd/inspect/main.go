package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/curusarn/resh/pkg/cfg"
	"github.com/curusarn/resh/pkg/msg"

	"os/user"
	"path/filepath"
	"strconv"
)

// version from git set during build
var version string

// commit from git set during build
var commit string

func main() {
	usr, _ := user.Current()
	dir := usr.HomeDir
	configPath := filepath.Join(dir, "/.config/resh.toml")

	var config cfg.Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Fatal("Error reading config:", err)
	}

	sessionID := flag.String("sessionID", "", "resh generated session id")
	count := flag.Uint("count", 10, "Number of cmdLines to return")
	flag.Parse()

	if *sessionID == "" {
		fmt.Println("Error: you need to specify sessionId")
	}

	m := msg.InspectMsg{SessionID: *sessionID, Count: *count}
	resp := SendInspectMsg(m, strconv.Itoa(config.Port))
	for _, cmdLine := range resp.CmdLines {
		fmt.Println("`" + cmdLine + "'")
	}
}

// SendInspectMsg to daemon
func SendInspectMsg(m msg.InspectMsg, port string) msg.MultiResponse {
	recJSON, err := json.Marshal(m)
	if err != nil {
		log.Fatal("send err 1", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:"+port+"/inspect",
		bytes.NewBuffer(recJSON))
	if err != nil {
		log.Fatal("send err 2", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("resh-daemon is not running :(")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("read response error")
	}
	// log.Println(string(body))
	response := msg.MultiResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal("unmarshal resp error: ", err)
	}
	return response
}
