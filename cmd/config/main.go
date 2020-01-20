package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/curusarn/resh/pkg/cfg"
)

func main() {
	usr, _ := user.Current()
	dir := usr.HomeDir
	configPath := filepath.Join(dir, ".config/resh.toml")

	var config cfg.Config
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		fmt.Println("Error reading config", err)
		os.Exit(1)
	}

	configKey := flag.String("key", "", "Key of the requested config entry")
	flag.Parse()

	if *configKey == "" {
		fmt.Println("Error: expected option --key!")
		os.Exit(1)
	}

	*configKey = strings.ToLower(*configKey)
	switch *configKey {
	case "bindarrowkeysbash":
		printBoolNormalized(config.BindArrowKeysBash)
	case "bindarrowkeyszsh":
		printBoolNormalized(config.BindArrowKeysZsh)
	case "bindcontrolr":
		printBoolNormalized(config.BindControlR)
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
