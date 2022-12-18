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
	"github.com/curusarn/resh/internal/datadir"
	"github.com/curusarn/resh/internal/device"
	"github.com/curusarn/resh/internal/logger"
	"github.com/curusarn/resh/internal/status"
	"go.uber.org/zap"
)

// info passed during build
var version string
var commit string
var development string

func main() {
	config, errCfg := cfg.New()
	logger, err := logger.New("daemon", config.LogLevel, development)
	if err != nil {
		fmt.Printf("Error while creating logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any
	if errCfg != nil {
		logger.Error("Error while getting configuration", zap.Error(errCfg))
	}
	sugar := logger.Sugar()
	d := daemon{sugar: sugar}
	sugar.Infow("Daemon starting ...",
		"version", version,
		"commit", commit,
	)
	dataDir, err := datadir.MakePath()
	if err != nil {
		sugar.Fatalw("Could not get user data directory", zap.Error(err))
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		sugar.Fatalw("Could not get user home directory", zap.Error(err))
	}
	// TODO: These paths should be probably defined in a package
	pidFile := filepath.Join(dataDir, "daemon.pid")
	reshHistoryPath := filepath.Join(dataDir, "history.reshjson")
	bashHistoryPath := filepath.Join(homeDir, ".bash_history")
	zshHistoryPath := filepath.Join(homeDir, ".zsh_history")
	deviceID, err := device.GetID(dataDir)
	if err != nil {
		sugar.Fatalw("Could not get resh device ID", zap.Error(err))
	}
	deviceName, err := device.GetName(dataDir)
	if err != nil {
		sugar.Fatalw("Could not get resh device name", zap.Error(err))
	}

	sugar = sugar.With(zap.Int("daemonPID", os.Getpid()))

	res, err := status.IsDaemonRunning(config.Port)
	if err != nil {
		sugar.Errorw("Error while checking daemon status - it's probably not running",
			"error", err)
	}
	if res {
		sugar.Errorw("Daemon is already running - exiting!")
		return
	}
	_, err = os.Stat(pidFile)
	if err == nil {
		sugar.Warnw("PID file exists",
			"PIDFile", pidFile)
		// kill daemon
		err = d.killDaemon(pidFile)
		if err != nil {
			sugar.Errorw("Could not kill daemon",
				"error", err,
			)
		}
	}
	err = ioutil.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0644)
	if err != nil {
		sugar.Fatalw("Could not create PID file",
			"error", err,
			"PIDFile", pidFile,
		)
	}
	server := Server{
		sugar:           sugar,
		config:          config,
		reshHistoryPath: reshHistoryPath,
		bashHistoryPath: bashHistoryPath,
		zshHistoryPath:  zshHistoryPath,

		deviceID:   deviceID,
		deviceName: deviceName,
	}
	server.Run()
	sugar.Infow("Removing PID file ...",
		"PIDFile", pidFile,
	)
	err = os.Remove(pidFile)
	if err != nil {
		sugar.Errorw("Could not delete PID file", "error", err)
	}
	sugar.Info("Shutting down ...")
}

type daemon struct {
	sugar *zap.SugaredLogger
}

func (d *daemon) killDaemon(pidFile string) error {
	dat, err := ioutil.ReadFile(pidFile)
	if err != nil {
		d.sugar.Errorw("Reading PID file failed",
			"PIDFile", pidFile,
			"error", err)
	}
	d.sugar.Infow("Successfully read PID file", "contents", string(dat))
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
