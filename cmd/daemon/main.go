package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/curusarn/resh/internal/cfg"
	"github.com/curusarn/resh/internal/httpclient"
	"github.com/curusarn/resh/internal/logger"
	"go.uber.org/zap"
)

// info passed during build
var version string
var commit string
var developement bool

func main() {
	config, errCfg := cfg.New()
	logger, err := logger.New("daemon", config.LogLevel, developement)
	if err != nil {
		fmt.Printf("Error while creating logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any
	if errCfg != nil {
		logger.Error("Error while getting configuration", zap.Error(errCfg))
	}
	sugar := logger.Sugar()
	d := daemon{sugar: sugar}
	sugar.Infow("Deamon starting ...",
		"version", version,
		"commit", commit,
	)

	// TODO: rethink PID file and logs location
	homeDir, err := os.UserHomeDir()
	if err != nil {
		sugar.Fatalw("Could not get user home dir", zap.Error(err))
	}
	PIDFile := filepath.Join(homeDir, ".resh/resh.pid")
	reshHistoryPath := filepath.Join(homeDir, ".resh_history.json")
	bashHistoryPath := filepath.Join(homeDir, ".bash_history")
	zshHistoryPath := filepath.Join(homeDir, ".zsh_history")

	sugar = sugar.With(zap.Int("daemonPID", os.Getpid()))

	res, err := d.isDaemonRunning(config.Port)
	if err != nil {
		sugar.Errorw("Error while checking daemon status - "+
			"it's probably not running", "error", err)
	}
	if res {
		sugar.Errorw("Daemon is already running - exiting!")
		return
	}
	_, err = os.Stat(PIDFile)
	if err == nil {
		sugar.Warn("Pidfile exists")
		// kill daemon
		err = d.killDaemon(PIDFile)
		if err != nil {
			sugar.Errorw("Could not kill daemon",
				"error", err,
			)
		}
	}
	err = ioutil.WriteFile(PIDFile, []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		sugar.Fatalw("Could not create pidfile",
			"error", err,
			"PIDFile", PIDFile,
		)
	}
	server := Server{
		sugar:           sugar,
		config:          config,
		reshHistoryPath: reshHistoryPath,
		bashHistoryPath: bashHistoryPath,
		zshHistoryPath:  zshHistoryPath,
	}
	server.Run()
	sugar.Infow("Removing PID file ...",
		"PIDFile", PIDFile,
	)
	err = os.Remove(PIDFile)
	if err != nil {
		sugar.Errorw("Could not delete PID file", "error", err)
	}
	sugar.Info("Shutting down ...")
}

type daemon struct {
	sugar *zap.SugaredLogger
}

func (d *daemon) getEnvOrPanic(envVar string) string {
	val, found := os.LookupEnv(envVar)
	if !found {
		d.sugar.Fatalw("Required env variable is not set",
			"variableName", envVar,
		)
	}
	return val
}

func (d *daemon) isDaemonRunning(port int) (bool, error) {
	url := "http://localhost:" + strconv.Itoa(port) + "/status"
	client := httpclient.New()
	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}

func (d *daemon) killDaemon(pidfile string) error {
	dat, err := ioutil.ReadFile(pidfile)
	if err != nil {
		d.sugar.Errorw("Reading pid file failed",
			"PIDFile", pidfile,
			"error", err)
	}
	d.sugar.Infow("Succesfully read PID file", "contents", string(dat))
	pid, err := strconv.Atoi(strings.TrimSuffix(string(dat), "\n"))
	if err != nil {
		return fmt.Errorf("could not parse PID file contents: %w", err)
	}
	d.sugar.Infow("Successfully parsed PID", "PID", pid)
	cmd := exec.Command("kill", "-s", "sigint", strconv.Itoa(pid))
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("kill command finished with error: %w", err)
	}
	return nil
}
