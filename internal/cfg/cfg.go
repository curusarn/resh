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
	// SesswatchPeriodSeconds is how often should daemon check if terminal
	// sessions are still alive
	SesswatchPeriodSeconds uint
	// SesshistInitHistorySize is how large resh history needs to be for
	// daemon to ignore standard shell history files
	SesshistInitHistorySize int
}

// defaults for config
var defaults = Config{
	Port:         2627,
	LogLevel:     zap.InfoLevel,
	BindControlR: true,

	Debug:                   false,
	SesswatchPeriodSeconds:  600,
	SesshistInitHistorySize: 1000,
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
