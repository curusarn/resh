package main

import (

	//"flag"
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
)

// Version from git set during build
var Version string

// Revision from git set during build
var Revision string

func main() {
	log.Println("Daemon starting... \n" +
		"version: " + Version +
		" revision: " + Revision)
	usr, _ := user.Current()
	dir := usr.HomeDir
	pidfilePath := filepath.Join(dir, ".resh/resh.pid")
	configPath := filepath.Join(dir, ".config/resh.toml")
	historyPath := filepath.Join(dir, ".resh_history.json")
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
	runServer(config, historyPath)
	log.Println("main: Removing pidfile ...")
	err = os.Remove(pidfilePath)
	if err != nil {
		log.Println("Could not delete pidfile", err)
	}
	log.Println("main: Shutdown - bye")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK; version: " + Version +
		"; revision: " + Revision + "\n"))
	log.Println("Status OK")
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
