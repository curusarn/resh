package cfg

import (
	"fmt"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// configFile used to parse the config file
type configFile struct {
	// ConfigVersion - never remove this
	ConfigVersion *string

	// added in legacy
	Port                    *int
	SesswatchPeriodSeconds  *uint
	SesshistInitHistorySize *int
	BindControlR            *bool
	Debug                   *bool

	// added in v1
	LogLevel *string

	// added in legacy
	// deprecated in v1
	BindArrowKeysBash *bool
	BindArrowKeysZsh  *bool

	SyncConnectorAddress *string
}

// Config returned by this package to be used in the rest of the project
type Config struct {
	// Port used by daemon and rest of the components to communicate
	// Make sure to restart the daemon when you change it
	Port int

	// BindControlR causes CTRL+R to launch the search app
	BindControlR bool
	// LogLevel used to filter logs
	LogLevel zapcore.Level

	// Debug mode for search app
	Debug bool
	// SessionWatchPeriodSeconds is how often should daemon check if terminal
	// sessions are still alive
	// There is not much need to adjust the value both memory overhead of watched sessions
	// and the CPU overhead of checking them are relatively low
	SessionWatchPeriodSeconds uint
	// ReshHistoryMinSize is how large resh history needs to be for
	// daemon to ignore standard shell history files
	// Ignoring standard shell history gives us more consistent experience,
	// but you can increase this to something large to see standard shell history in RESH search
	ReshHistoryMinSize int

	// SyncConnectorAddress used by the daemon to connect to the Sync Connector
	// examples:
	//  - localhost:1234
	//  - http://localhost:1234
	//  - 192.168.1.1:1324
	//  - https://domain.tld
	//  - https://domain.tld/resh
	SyncConnectorAddress *string
}

// defaults for config
var defaults = Config{
	Port:         2627,
	LogLevel:     zap.InfoLevel,
	BindControlR: true,

	Debug:                     false,
	SessionWatchPeriodSeconds: 600,
	ReshHistoryMinSize:        1000,
}

const headerComment = `##
######################
## RESH config (v1) ##
######################
## Here you can find info about RESH configuration options.
## You can uncomment the options and custimize them.

## Required.
## The config format can change in future versions.
## ConfigVersion helps us seemlessly upgrade to the new formats.
# ConfigVersion = "v1"

## Port used by RESH daemon and rest of the components to communicate.
## Make sure to restart the daemon (pkill resh-daemon) when you change it.
# Port = 2627

## Controls how much and how detailed logs all RESH components produce.
## Use "debug" for full logs when you encounter an issue
## Options: "debug", "info", "warn", "error", "fatal"
# LogLevel = "info"

## When BindControlR is "true" RESH search app is bound to CTRL+R on terminal startuA
# BindControlR = true

## When Debug is "true" the RESH search app runs in debug mode.
## This is useful for development.
# Debug = false

## Daemon keeps track of running terminal sessions.
## SessionWatchPeriodSeconds controls how often daemon checks if the sessions are still alive.
## You shouldn't need to adjust this.
# SessionWatchPeriodSeconds = 600

## When RESH is first installed there is no RESH history so there is nothing to search.
## As a temporary woraround, RESH daemon parses bash/zsh shell history and searches it.
## Once RESH history is big enough RESH stops using bash/zsh history.
## ReshHistoryMinSize controls how big RESH history needs to be before this happens.
## You can increase this this to e.g. 10000 to get RESH to use bash/zsh history longer.
# ReshHistoryMinSize = 1000

`

func getConfigPath() (string, error) {
	fname := "resh.toml"
	xdgDir, found := os.LookupEnv("XDG_CONFIG_HOME")
	if found {
		return path.Join(xdgDir, fname), nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home dir: %w", err)
	}
	return path.Join(homeDir, ".config", fname), nil
}

func readConfig(path string) (*configFile, error) {
	var config configFile
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return &config, fmt.Errorf("could not decode config: %w", err)
	}
	return &config, nil
}

func getConfig() (*configFile, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("could not get config file path: %w", err)
	}
	return readConfig(path)
}

// returned config is always usable, returned errors are informative
func processAndFillDefaults(configF *configFile) (Config, error) {
	config := defaults

	if configF.Port != nil {
		config.Port = *configF.Port
	}
	if configF.SesswatchPeriodSeconds != nil {
		config.SessionWatchPeriodSeconds = *configF.SesswatchPeriodSeconds
	}
	if configF.SesshistInitHistorySize != nil {
		config.ReshHistoryMinSize = *configF.SesshistInitHistorySize
	}

	var err error
	if configF.LogLevel != nil {
		logLevel, err := zapcore.ParseLevel(*configF.LogLevel)
		if err != nil {
			err = fmt.Errorf("could not parse log level: %w", err)
		} else {
			config.LogLevel = logLevel
		}
	}

	if configF.BindControlR != nil {
		config.BindControlR = *configF.BindControlR
	}

	config.SyncConnectorAddress = configF.SyncConnectorAddress

	return config, err
}

// New returns a config file
// returned config is always usable, returned errors are informative
func New() (Config, error) {
	configF, err := getConfig()
	if err != nil {
		return defaults, fmt.Errorf("using default config because of error while getting/reading config: %w", err)
	}

	config, err := processAndFillDefaults(configF)
	if err != nil {
		return config, fmt.Errorf("errors while processing config: %w", err)
	}
	return config, nil
}
