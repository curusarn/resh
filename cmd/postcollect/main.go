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

	"strconv"
)

// info passed during build
var version string
var commit string
var developement bool

func main() {
	config, errCfg := cfg.New()
	logger, _ := logger.New("postcollect", config.LogLevel, developement)
	defer logger.Sync() // flushes buffer, if any
	if errCfg != nil {
		logger.Error("Error while getting configuration", zap.Error(errCfg))
	}
	out := output.New(logger, "resh-postcollect ERROR")

	showVersion := flag.Bool("version", false, "Show version and exit")
	showRevision := flag.Bool("revision", false, "Show git revision and exit")

	requireVersion := flag.String("requireVersion", "", "abort if version doesn't match")
	requireRevision := flag.String("requireRevision", "", "abort if revision doesn't match")

	exitCode := flag.Int("exitCode", -1, "exit code")
	sessionID := flag.String("sessionID", "", "resh generated session id")
	recordID := flag.String("recordID", "", "resh generated record id")

	shlvl := flag.Int("shlvl", -1, "$SHLVL")

	rtb := flag.String("timeBefore", "-1", "before $EPOCHREALTIME")
	rta := flag.String("timeAfter", "-1", "after $EPOCHREALTIME")
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

	timeAfter, err := strconv.ParseFloat(*rta, 64)
	if err != nil {
		out.Fatal("Error while parsing flag --timeAfter", err)
	}
	timeBefore, err := strconv.ParseFloat(*rtb, 64)
	if err != nil {
		out.Fatal("Error while parsing flag --timeBefore", err)
	}
	duration := timeAfter - timeBefore

	// FIXME: use recordint.Postollect
	rec := recordint.Collect{
		SessionID: *sessionID,
		Shlvl:     *shlvl,

		Rec: record.V1{
			RecordID:  *recordID,
			SessionID: *sessionID,

			ExitCode: *exitCode,
			Duration: duration,

			PartsNotMerged: true,
		},
	}
	collect.SendRecord(out, rec, strconv.Itoa(config.Port), "/record")
}
