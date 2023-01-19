package main

import (
	"fmt"
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

const helpMsg = `ERROR: resh-daemon doesn't accept any arguments

WARNING:
  You shouldn't typically need to start RESH daemon yourself.
  Unless its already running, RESH daemon is started when a new terminal is opened.
  RESH daemon will not start if it's already running even when you run it manually.

USAGE:
  $ resh-daemon
  Runs the daemon as foreground process. You can kill it with CTRL+C.

  $ resh-daemon-start
  Runs the daemon as background process detached from terminal.

LOGS & DEBUGGING:
  Logs are located in:
    ${XDG_DATA_HOME}/resh/log.json (if XDG_DATA_HOME is set)
    ~/.local/share/resh/log.json   (otherwise - more common)

  A good way to see the logs as they are being produced is:
    $ tail -f ~/.local/share/resh/log.json

MORE INFO:
  https://github.com/curusarn/resh/
`

func main() {
	if len(os.Args) > 1 {
		fmt.Fprint(os.Stderr, helpMsg)
		os.Exit(1)
	}
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
	reshHistoryPath := filepath.Join(dataDir, datadir.HistoryFileName)
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
	err = os.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0644)
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
	dat, err := os.ReadFile(pidFile)
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
	err = exec.Command("kill", "-SIGTERM", fmt.Sprintf("%d", pid)).Run()
	if err != nil {
		return fmt.Errorf("kill command finished with error: %w", err)
	}
	return nil
}
