package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

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
