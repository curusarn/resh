package main

import (
	"fmt"
	"os"

	"github.com/curusarn/resh/internal/cfg"
	"github.com/curusarn/resh/internal/collect"
	"github.com/curusarn/resh/internal/logger"
	"github.com/curusarn/resh/internal/opt"
	"github.com/curusarn/resh/internal/output"
	"github.com/curusarn/resh/internal/recordint"
	"github.com/curusarn/resh/record"
	"github.com/spf13/pflag"
	"go.uber.org/zap"

	"path/filepath"
	"strconv"
)

// info passed during build
var version string
var commit string
var development string

func main() {
	config, errCfg := cfg.New()
	logger, err := logger.New("collect", config.LogLevel, development)
	if err != nil {
		fmt.Printf("Error while creating logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any
	if errCfg != nil {
		logger.Error("Error while getting configuration", zap.Error(errCfg))
	}
	out := output.New(logger, "resh-collect ERROR")

	args := opt.HandleVersionOpts(out, os.Args, version, commit)

	flags := pflag.NewFlagSet("", pflag.ExitOnError)
	cmdLine := flags.String("cmd-line", "", "Command line")
	gitRemote := flags.String("git-remote", "", "> git remote get-url origin")
	home := flags.String("home", "", "$HOME")
	pwd := flags.String("pwd", "", "$PWD - present working directory")
	recordID := flags.String("record-id", "", "Resh generated record ID")
	sessionID := flags.String("session-id", "", "Resh generated session ID")
	sessionPID := flags.Int("session-pid", -1, "$$ - Shell session PID")
	shell := flags.String("shell", "", "Current shell")
	shlvl := flags.Int("shlvl", -1, "$SHLVL")
	timeStr := flags.String("time", "-1", "$EPOCHREALTIME")
	flags.Parse(args)

	time, err := strconv.ParseFloat(*timeStr, 64)
	if err != nil {
		out.FatalE("Error while parsing flag --time", err)
	}

	realPwd, err := filepath.EvalSymlinks(*pwd)
	if err != nil {
		out.ErrorE("Error while evaluating symlinks in PWD", err)
		realPwd = ""
	}

	rec := recordint.Collect{
		SessionID:  *sessionID,
		Shlvl:      *shlvl,
		SessionPID: *sessionPID,

		Shell: *shell,

		Rec: record.V1{
			SessionID: *sessionID,
			RecordID:  *recordID,

			CmdLine: *cmdLine,

			// posix
			Home:    *home,
			Pwd:     *pwd,
			RealPwd: realPwd,

			GitOriginRemote: *gitRemote,

			Time: fmt.Sprintf("%.4f", time),

			PartOne:        true,
			PartsNotMerged: true,
		},
	}
	collect.SendRecord(out, rec, strconv.Itoa(config.Port), "/record")
}
