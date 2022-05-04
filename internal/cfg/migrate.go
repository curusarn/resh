package cfg

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Touch config file
func Touch() error {
	path, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("could not get config file path: %w", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("could not open/create config file: %w", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("could not close config file: %w", err)
	}
	return nil
}

// Migrate old config versions to current config version
// returns true if any changes were made to the config
func Migrate() (bool, error) {
	path, err := getConfigPath()
	if err != nil {
		return false, fmt.Errorf("could not get config file path: %w", err)
	}
	configF, err := readConfig(path)
	if err != nil {
		return false, fmt.Errorf("could not read config: %w", err)
	}
	const current = "v1"
	if configF.ConfigVersion != nil && *configF.ConfigVersion == current {
		return false, nil
	}

	if configF.ConfigVersion == nil {
		configF, err = legacyToV1(configF)
		if err != nil {
			return true, fmt.Errorf("error converting config from version 'legacy' to 'v1': %w", err)
		}
	}

	if *configF.ConfigVersion != current {
		return false, fmt.Errorf("unrecognized config version: '%s'", *configF.ConfigVersion)
	}
	err = writeConfig(configF, path)
	if err != nil {
		return true, fmt.Errorf("could not write migrated config: %w", err)
	}
	return true, nil
}

func writeConfig(config *configFile, path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("could not open config for writing: %w", err)
	}
	defer file.Close()
	err = toml.NewEncoder(file).Encode(config)
	if err != nil {
		return fmt.Errorf("could not encode config: %w", err)
	}
	return nil
}

func legacyToV1(config *configFile) (*configFile, error) {
	if config.ConfigVersion != nil {
		return nil, fmt.Errorf("config version is not 'legacy': '%s'", *config.ConfigVersion)
	}
	version := "v1"
	newConf := configFile{
		ConfigVersion: &version,
	}
	// Remove defaults
	if config.Port != nil && *config.Port != 2627 {
		newConf.Port = config.Port
	}
	if config.SesswatchPeriodSeconds != nil && *config.SesswatchPeriodSeconds != 120 {
		newConf.SesswatchPeriodSeconds = config.SesswatchPeriodSeconds
	}
	if config.SesshistInitHistorySize != nil && *config.SesshistInitHistorySize != 1000 {
		newConf.SesshistInitHistorySize = config.SesshistInitHistorySize
	}
	if config.BindControlR != nil && *config.BindControlR != true {
		newConf.BindControlR = config.BindControlR
	}
	if config.Debug != nil && *config.Debug != false {
		newConf.Debug = config.Debug
	}
	return &newConf, nil
}

// func v1ToV2(config *configFile) (*configFile, error) {
// 	if *config.ConfigVersion != "v1" {
// 		return nil, fmt.Errorf("config version is not 'legacy': '%s'", *config.ConfigVersion)
// 	}
// 	version := "v2"
// 	newConf := configFile{
// 		ConfigVersion: &version,
// 		// Here goes all config fields - no need to prune defaults like we do for legacy
// 	}
// 	return &newConf, nil
// }
