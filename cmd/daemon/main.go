package main

import (
	//"flag"

	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/curusarn/resh/pkg/cfg"
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
