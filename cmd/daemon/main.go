package main

import (
	"encoding/json"
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
	"github.com/curusarn/resh/pkg/histfile"
	"github.com/curusarn/resh/pkg/records"
	"github.com/curusarn/resh/pkg/sesswatch"
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
	outputPath := filepath.Join(dir, ".resh_history.json")
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
	runServer(config, outputPath)
	err = os.Remove(pidfilePath)
	if err != nil {
		log.Println("Could not delete pidfile", err)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK; version: " + Version +
		"; revision: " + Revision + "\n"))
	log.Println("Status OK")
}

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
		log.Println("Received:", record.CmdLine, " - part", part)
	}()

	// fmt.Println("cmd:", r.CmdLine)
	// fmt.Println("pwd:", r.Pwd)
	// fmt.Println("git:", r.GitWorkTree)
	// fmt.Println("exit_code:", r.ExitCode)
}

func runServer(config cfg.Config, outputPath string) {
	var recordSubscribers []chan records.Record

	histfileChan := make(chan records.Record)
	recordSubscribers = append(recordSubscribers, histfileChan)
	sessionsToDrop := make(chan string)
	histfile.Go(histfileChan, outputPath, sessionsToDrop)

	sesswatchChan := make(chan records.Record)
	recordSubscribers = append(recordSubscribers, sesswatchChan)
	sesswatch.Go(sesswatchChan, []chan string{sessionsToDrop}, config.SesswatchPeriodSeconds)

	http.HandleFunc("/status", statusHandler)
	http.Handle("/record", &recordHandler{subscribers: recordSubscribers})
	//http.Handle("/session_init", &sessionInitHandler{OutputPath: outputPath})
	//http.Handle("/recall", &recallHandler{OutputPath: outputPath})
	http.ListenAndServe(":"+strconv.Itoa(config.Port), nil)
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
	cmd := exec.Command("kill", strconv.Itoa(pid))
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
	//body, err := ioutil.ReadAll(resp.Body)

	//    dat, err := ioutil.ReadFile(pidfile)
	//    if err != nil {
	//        log.Println("Reading pid file failed", err)
	//        return false, err
	//    }
	//    log.Print(string(dat))
	//    pid, err := strconv.ParseInt(string(dat), 10, 64)
	//    if err != nil {
	//        log.Fatal(err)
	//    }
	//    process, err := os.FindProcess(int(pid))
	//    if err != nil {
	//        log.Printf("Failed to find process: %s\n", err)
	//        return false, err
	//    } else {
	//        err := process.Signal(syscall.Signal(0))
	//        log.Printf("process.Signal on pid %d returned: %v\n", pid, err)
	//    }
	//    return true, nil
}
