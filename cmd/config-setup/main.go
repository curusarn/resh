package main

import (
	"fmt"
	"os"

	"github.com/curusarn/resh/internal/cfg"
	"github.com/curusarn/resh/internal/logger"
	"go.uber.org/zap"
)

// info passed during build
var version string
var commit string
var developement bool

func main() {
	errDo := doConfigSetup()
	config, errCfg := cfg.New()
	logger, _ := logger.New("config-setup", config.LogLevel, developement)
	defer logger.Sync() // flushes buffer, if any

	if errDo != nil {
		logger.Error("Config setup failed", zap.Error(errDo))
		// TODO: better error message for people
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", errDo)
	}
	if errCfg != nil {
		logger.Error("Error while getting configuration", zap.Error(errCfg))
	}
}

func doConfigSetup() error {
	err := cfg.Touch()
	if err != nil {
		return fmt.Errorf("could not touch config file: %w", err)
	}
	changes, err := cfg.Migrate()
	if err != nil {
		return fmt.Errorf("could not migrate config file version: %v", err)
	}
	if changes {
		fmt.Printf("Config file format has changed - your config was updated to reflect the changes.\n")
	}
	return nil
}
