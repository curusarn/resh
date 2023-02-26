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

	"strconv"
)

// info passed during build
var version string
var commit string
var development string

func main() {
	config, errCfg := cfg.New()
	logger, err := logger.New("postcollect", config.LogLevel, development)
	if err != nil {
		fmt.Printf("Error while creating logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any
	if errCfg != nil {
		logger.Error("Error while getting configuration", zap.Error(errCfg))
	}
	out := output.New(logger, "resh-postcollect ERROR")

	args := opt.HandleVersionOpts(out, os.Args, version, commit)

	flags := pflag.NewFlagSet("", pflag.ExitOnError)
	exitCode := flags.Int("exit-code", -1, "Exit code")
	sessionID := flags.String("session-id", "", "Resh generated session ID")
	recordID := flags.String("record-id", "", "Resh generated record ID")
	shlvl := flags.Int("shlvl", -1, "$SHLVL")
	rtb := flags.String("time-before", "-1", "Before $EPOCHREALTIME")
	rta := flags.String("time-after", "-1", "After $EPOCHREALTIME")
	flags.Parse(args)

	timeAfter, err := strconv.ParseFloat(*rta, 64)
	if err != nil {
		out.FatalE("Error while parsing flag --time-after", err)
	}
	timeBefore, err := strconv.ParseFloat(*rtb, 64)
	if err != nil {
		out.FatalE("Error while parsing flag --time-before", err)
	}
	duration := timeAfter - timeBefore

	// FIXME: use recordint.Postcollect
	rec := recordint.Collect{
		SessionID: *sessionID,
		Shlvl:     *shlvl,

		Rec: record.V1{
			RecordID:  *recordID,
			SessionID: *sessionID,

			ExitCode: *exitCode,
			Duration: fmt.Sprintf("%.4f", duration),

			PartsNotMerged: true,
		},
	}
	collect.SendRecord(out, rec, strconv.Itoa(config.Port), "/record")
}
