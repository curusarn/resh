package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/curusarn/resh/cmd/control/status"
	"github.com/curusarn/resh/pkg/cfg"
	"github.com/spf13/cobra"
)

// Enable commands

var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "enable RESH features (bindings)",
}

var enableControlRBindingCmd = &cobra.Command{
	Use:   "ctrl_r_binding",
	Short: "enable RESH-CLI binding for Ctrl+R",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = enableDisableControlRBindingGlobally(true)
		if exitCode == status.Success {
			exitCode = status.EnableControlRBinding
		}
	},
}

// Disable commands

var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "disable RESH features (bindings)",
}

var disableControlRBindingCmd = &cobra.Command{
	Use:   "ctrl_r_binding",
	Short: "disable RESH-CLI binding for Ctrl+R",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = enableDisableControlRBindingGlobally(false)
		if exitCode == status.Success {
			exitCode = status.DisableControlRBinding
		}
	},
}

func enableDisableControlRBindingGlobally(value bool) status.Code {
	usr, _ := user.Current()
	dir := usr.HomeDir
	configPath := filepath.Join(dir, ".config/resh.toml")
	var config cfg.Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		fmt.Println("Error reading config", err)
		return status.Fail
	}
	if config.BindControlR != value {
		config.BindControlR = value

		f, err := os.Create(configPath)
		if err != nil {
			fmt.Println("Error: Failed to create/open file:", configPath, "; error:", err)
			return status.Fail
		}
		defer f.Close()
		if err := toml.NewEncoder(f).Encode(config); err != nil {
			fmt.Println("Error: Failed to encode and write the config values to hdd. error:", err)
			return status.Fail
		}
	}
	if value {
		fmt.Println("RESH SEARCH app Ctrl+R binding: ENABLED")
	} else {
		fmt.Println("RESH SEARCH app Ctrl+R binding: DISABLED")
	}
	return status.Success
}
