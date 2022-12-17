package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/curusarn/resh/internal/cfg"
	"github.com/curusarn/resh/internal/collect"
	"github.com/curusarn/resh/internal/logger"
	"github.com/curusarn/resh/internal/output"
	"github.com/curusarn/resh/internal/recordint"
	"go.uber.org/zap"

	"strconv"
)

// info passed during build
var version string
var commit string
var development string

func main() {
	config, errCfg := cfg.New()
	logger, _ := logger.New("collect", config.LogLevel, development)
	defer logger.Sync() // flushes buffer, if any
	if errCfg != nil {
		logger.Error("Error while getting configuration", zap.Error(errCfg))
	}
	out := output.New(logger, "resh-collect ERROR")

	showVersion := flag.Bool("version", false, "Show version and exit")
	showRevision := flag.Bool("revision", false, "Show git revision and exit")

	requireVersion := flag.String("requireVersion", "", "abort if version doesn't match")
	requireRevision := flag.String("requireRevision", "", "abort if revision doesn't match")

	sessionID := flag.String("sessionId", "", "resh generated session id")

	sessionPID := flag.Int("sessionPid", -1, "$$ at session start")
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
		out.FatalTerminalVersionMismatch(version, *requireVersion)
	}
	if *requireRevision != "" && *requireRevision != commit {
		// this is only relevant for dev versions so we can reuse FatalVersionMismatch()
		out.FatalTerminalVersionMismatch("revision "+commit, "revision "+*requireVersion)
	}

	rec := recordint.SessionInit{
		SessionID:  *sessionID,
		SessionPID: *sessionPID,
	}
	collect.SendSessionInit(out, rec, strconv.Itoa(config.Port))
}
