package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/curusarn/resh/internal/cfg"
	"github.com/curusarn/resh/internal/collect"
	"github.com/curusarn/resh/internal/logger"
	"github.com/curusarn/resh/internal/output"
	"github.com/curusarn/resh/internal/record"
	"github.com/curusarn/resh/internal/recordint"
	"go.uber.org/zap"

	//  "os/exec"

	"path/filepath"
	"strconv"
)

// info passed during build
var version string
var commit string
var developement bool

func main() {
	config, errCfg := cfg.New()
	logger, _ := logger.New("collect", config.LogLevel, developement)
	defer logger.Sync() // flushes buffer, if any
	if errCfg != nil {
		logger.Error("Error while getting configuration", zap.Error(errCfg))
	}
	out := output.New(logger, "resh-collect ERROR")

	// version
	showVersion := flag.Bool("version", false, "Show version and exit")
	showRevision := flag.Bool("revision", false, "Show git revision and exit")

	requireVersion := flag.String("requireVersion", "", "abort if version doesn't match")
	requireRevision := flag.String("requireRevision", "", "abort if revision doesn't match")

	// core
	cmdLine := flag.String("cmdLine", "", "command line")

	home := flag.String("home", "", "$HOME")
	pwd := flag.String("pwd", "", "$PWD - present working directory")

	// FIXME: get device ID
	deviceID := flag.String("deviceID", "", "RESH device ID")
	sessionID := flag.String("sessionID", "", "resh generated session ID")
	recordID := flag.String("recordID", "", "resh generated record ID")
	sessionPID := flag.Int("sessionPID", -1, "PID at the start of the terminal session")

	shell := flag.String("shell", "", "current shell")

	// logname := flag.String("logname", "", "$LOGNAME")
	device := flag.String("device", "", "device name, usually $HOSTNAME")

	// non-posix
	shlvl := flag.Int("shlvl", -1, "$SHLVL")

	gitRemote := flag.String("gitRemote", "", "git remote get-url origin")

	time_ := flag.String("time", "-1", "$EPOCHREALTIME")
	flag.Parse()

	if *showVersion == true {
		fmt.Println(version)
		os.Exit(0)
	}
	if *showRevision == true {
		fmt.Println(commit)
		os.Exit(0)
	}
	if *requireVersion != "" && *requireVersion != version {
		out.FatalVersionMismatch(version, *requireVersion)
	}
	if *requireRevision != "" && *requireRevision != commit {
		// this is only relevant for dev versions so we can reuse FatalVersionMismatch()
		out.FatalVersionMismatch("revision "+commit, "revision "+*requireVersion)
	}

	time, err := strconv.ParseFloat(*time_, 64)
	if err != nil {
		out.Fatal("Error while parsing flag --time", err)
	}

	realPwd, err := filepath.EvalSymlinks(*pwd)
	if err != nil {
		logger.Error("Error while handling pwd realpath", zap.Error(err))
		realPwd = ""
	}

	rec := recordint.Collect{
		SessionID:  *sessionID,
		Shlvl:      *shlvl,
		SessionPID: *sessionPID,

		Shell: *shell,

		Rec: record.V1{
			DeviceID:  *deviceID,
			SessionID: *sessionID,
			RecordID:  *recordID,

			CmdLine: *cmdLine,

			// posix
			Home:    *home,
			Pwd:     *pwd,
			RealPwd: realPwd,

			// Logname:  *logname,
			Device: *device,

			GitOriginRemote: *gitRemote,

			Time: fmt.Sprintf("%.4f", time),

			PartOne:        true,
			PartsNotMerged: true,
		},
	}
	collect.SendRecord(out, rec, strconv.Itoa(config.Port), "/record")
}
