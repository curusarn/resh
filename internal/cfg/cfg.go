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
	Port                    *int
	SesswatchPeriodSeconds  *uint
	SesshistInitHistorySize *int

	LogLevel *string

	BindControlR *bool

	// deprecated
	BindArrowKeysBash *bool
	BindArrowKeysZsh  *bool
	Debug             *bool
}

// Config returned by this package to be used in the rest of the project
type Config struct {
	Port                    int
	SesswatchPeriodSeconds  uint
	SesshistInitHistorySize int
	LogLevel                zapcore.Level
	Debug                   bool
	BindControlR            bool
}

// defaults for config
var defaults = Config{
	Port:                    2627,
	SesswatchPeriodSeconds:  120,
	SesshistInitHistorySize: 1000,
	LogLevel:                zap.InfoLevel,
	Debug:                   false,
	BindControlR:            true,
}

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

func readConfig() (*configFile, error) {
	var config configFile
	path, err := getConfigPath()
	if err != nil {
		return &config, fmt.Errorf("could not get config file path: %w", err)
	}
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return &config, fmt.Errorf("could not decode config: %w", err)
	}
	return &config, nil
}

func processAndFillDefaults(configF *configFile) (Config, error) {
	config := defaults

	if configF.Port != nil {
		config.Port = *configF.Port
	}
	if configF.SesswatchPeriodSeconds != nil {
		config.SesswatchPeriodSeconds = *configF.SesswatchPeriodSeconds
	}
	if configF.SesshistInitHistorySize != nil {
		config.SesshistInitHistorySize = *configF.SesshistInitHistorySize
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

	return config, err
}

// New returns a config file
// returned config is always usable, returned errors are informative
func New() (Config, error) {
	configF, err := readConfig()
	if err != nil {
		return defaults, fmt.Errorf("using default config because of error while getting config: %w", err)
	}

	config, err := processAndFillDefaults(configF)
	if err != nil {
		return config, fmt.Errorf("errors while processing config: %w", err)
	}
	return config, nil
}
