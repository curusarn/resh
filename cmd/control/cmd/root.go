package cmd

import (
	"fmt"
	"log"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/curusarn/resh/cmd/control/status"
	"github.com/curusarn/resh/pkg/cfg"
	"github.com/spf13/cobra"
)

// globals
var exitCode status.Code
var version string
var commit string
var debug = false
var config cfg.Config

var rootCmd = &cobra.Command{
	Use:   "reshctl",
	Short: "Reshctl (RESH control) - check status, update, enable/disable features, sanitize history and more.",
}

// Execute reshctl
func Execute(ver, com string) status.Code {
	version = ver
	commit = com

	usr, _ := user.Current()
	dir := usr.HomeDir
	configPath := filepath.Join(dir, ".config/resh.toml")
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Println("Error reading config", err)
		return status.Fail
	}
	if config.Debug {
		debug = true
		// log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	}

	rootCmd.AddCommand(enableCmd)
	enableCmd.AddCommand(enableArrowKeyBindingsCmd)
	enableCmd.AddCommand(enableArrowKeyBindingsGlobalCmd)
	enableCmd.AddCommand(enableControlRBindingCmd)
	enableCmd.AddCommand(enableControlRBindingGlobalCmd)

	rootCmd.AddCommand(disableCmd)
	disableCmd.AddCommand(disableArrowKeyBindingsCmd)
	disableCmd.AddCommand(disableArrowKeyBindingsGlobalCmd)
	disableCmd.AddCommand(disableControlRBindingCmd)
	disableCmd.AddCommand(disableControlRBindingGlobalCmd)

	rootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(completionBashCmd)
	completionCmd.AddCommand(completionZshCmd)

	rootCmd.AddCommand(debugCmd)
	debugCmd.AddCommand(debugReloadCmd)
	debugCmd.AddCommand(debugInspectCmd)
	debugCmd.AddCommand(debugOutputCmd)

	rootCmd.AddCommand(statusCmd)

	rootCmd.AddCommand(updateCmd)

	rootCmd.AddCommand(sanitizeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return status.Fail
	}
	return exitCode
}
