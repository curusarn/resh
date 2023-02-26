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
	logger, err := logger.New("session-init", config.LogLevel, development)
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
	sessionID := flags.String("session-id", "", "RESH generated session ID")
	sessionPID := flags.Int("session-pid", -1, "$$ - Shell session PID")
	flags.Parse(args)

	rec := recordint.SessionInit{
		SessionID:  *sessionID,
		SessionPID: *sessionPID,
	}
	collect.SendSessionInit(out, rec, strconv.Itoa(config.Port))
}
