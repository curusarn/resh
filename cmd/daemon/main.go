package main

import (

	//"flag"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/curusarn/resh/pkg/cfg"
	"github.com/curusarn/resh/pkg/msg"
)

// version from git set during build
var version string

// commit from git set during build
var commit string

// Debug switch
var Debug = false

func main() {
	log.Println("Daemon starting... \n" +
		"version: " + version +
		" commit: " + commit)
	usr, _ := user.Current()
	dir := usr.HomeDir
	pidfilePath := filepath.Join(dir, ".resh/resh.pid")
	configPath := filepath.Join(dir, ".config/resh.toml")
	reshHistoryPath := filepath.Join(dir, ".resh_history.json")
	bashHistoryPath := filepath.Join(dir, ".bash_history")
	zshHistoryPath := filepath.Join(dir, ".zsh_history")
	logPath := filepath.Join(dir, ".resh/daemon.log")

	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.SetPrefix(strconv.Itoa(os.Getpid()) + " | ")

	var config cfg.Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Println("Error reading config", err)
		return
	}
	if config.Debug {
		Debug = true
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	}

	res, err := isDaemonRunning(config.Port)
	if err != nil {
		log.Println("Error while checking if the daemon is runnnig", err)
	}
	if res {
		log.Println("Daemon is already running - exiting!")
		return
	}
	_, err = os.Stat(pidfilePath)
	if err == nil {
		log.Println("Pidfile exists")
		// kill daemon
		err = killDaemon(pidfilePath)
		if err != nil {
			log.Println("Error while killing daemon", err)
		}
	}
	err = ioutil.WriteFile(pidfilePath, []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		log.Fatal("Could not create pidfile", err)
	}
	runServer(config, reshHistoryPath, bashHistoryPath, zshHistoryPath)
	log.Println("main: Removing pidfile ...")
	err = os.Remove(pidfilePath)
	if err != nil {
		log.Println("Could not delete pidfile", err)
	}
	log.Println("main: Shutdown - bye")
}

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

func killDaemon(pidfile string) error {
	dat, err := ioutil.ReadFile(pidfile)
	if err != nil {
		log.Println("Reading pid file failed", err)
	}
	log.Print(string(dat))
	pid, err := strconv.Atoi(strings.TrimSuffix(string(dat), "\n"))
	if err != nil {
		log.Fatal("Pidfile contents are malformed", err)
	}
	cmd := exec.Command("kill", "-s", "sigint", strconv.Itoa(pid))
	err = cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
		return err
	}
	return nil
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
