package main

import (
	"fmt"
	"os"

	"github.com/curusarn/resh/internal/cfg"
	"github.com/curusarn/resh/internal/logger"
	"github.com/curusarn/resh/internal/output"
	"go.uber.org/zap"
)

// info passed during build
var version string
var commit string
var development string

func main() {
	config, errCfg := cfg.New()
	logger, err := logger.New("install-utils", config.LogLevel, development)
	if err != nil {
		fmt.Printf("Error while creating logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any
	if errCfg != nil {
		logger.Error("Error while getting configuration", zap.Error(errCfg))
	}
	sugar := logger.Sugar()
	sugar.Infow("Install-utils invoked ...",
		"version", version,
		"commit", commit,
	)
	out := output.New(logger, "install-utils ERROR")

	if len(os.Args) < 2 {
		out.Error("ERROR: Not enough arguments\n")
		printUsage(os.Stderr)
		os.Exit(1)
	}
	command := os.Args[1]
	switch command {
	case "setup-device":
		setupDevice(out)
	case "migrate-all":
		migrateAll(out)
	case "help":
		printUsage(os.Stdout)
	default:
		out.Error(fmt.Sprintf("ERROR: Unknown command: %s\n", command))
		printUsage(os.Stderr)
		os.Exit(1)
	}
}

func printUsage(f *os.File) {
	usage := `
USAGE: ./install-utils COMMAND
Utils used during RESH installation.

COMMANDS:
  setup-device      setup device name and device ID
  migrate-all       update config and history to latest format
  help              show this help

`
	fmt.Fprint(f, usage)
}
