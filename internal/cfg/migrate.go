package cfg

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Migrate old config versions to current config version
// returns true if any changes were made to the config
func Migrate() (bool, error) {
	fpath, err := getConfigPath()
	if err != nil {
		return false, fmt.Errorf("could not get config file path: %w", err)
	}
	configF, err := readConfig(fpath)
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
	err = writeConfig(configF, fpath)
	if err != nil {
		return true, fmt.Errorf("could not write migrated config: %w", err)
	}
	return true, nil
}

// writeConfig should only be used when migrating config to new version
// writing the config file discards all comments in the config file (limitation of TOML library)
// to make up for lost comments we add header comment to the start of the file
func writeConfig(config *configFile, path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("could not open config for writing: %w", err)
	}
	defer file.Close()
	_, err = file.WriteString(headerComment)
	if err != nil {
		return fmt.Errorf("could not write config header: %w", err)
	}
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
