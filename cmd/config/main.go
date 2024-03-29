package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/curusarn/resh/internal/cfg"
	"github.com/curusarn/resh/internal/logger"
	"go.uber.org/zap"
)

// info passed during build
var version string
var commit string
var development string

func main() {
	config, errCfg := cfg.New()
	logger, _ := logger.New("config", config.LogLevel, development)
	defer logger.Sync() // flushes buffer, if any
	if errCfg != nil {
		logger.Error("Error while getting configuration", zap.Error(errCfg))
	}

	configKey := flag.String("key", "", "Key of the requested config entry")
	flag.Parse()

	if *configKey == "" {
		fmt.Println("Error: expected option --key!")
		os.Exit(1)
	}

	*configKey = strings.ToLower(*configKey)
	switch *configKey {
	case "bindcontrolr":
		printBoolNormalized(config.BindControlR)
	case "port":
		fmt.Println(config.Port)
	case "sesswatchperiodseconds":
		fmt.Println(config.SessionWatchPeriodSeconds)
	case "sesshistinithistorysize":
		fmt.Println(config.ReshHistoryMinSize)
	default:
		fmt.Println("Error: illegal --key!")
		os.Exit(1)
	}
}

// this might be unnecessary but I'm too lazy to look it up
func printBoolNormalized(x bool) {
	if x {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}
