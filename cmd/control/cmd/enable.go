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

var enableArrowKeyBindingsCmd = &cobra.Command{
	Use:   "arrow_key_bindings",
	Short: "enable bindings for arrow keys (up/down)",
	Long: "Enable bindings for arrow keys (up/down)\n" +
		"Note that this only affects sessions of the same shell.\n" +
		"(e.g. running this in zsh will only enable the keybinding in zsh)",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = enableDisableArrowKeyBindingsGlobally(true)
		if exitCode == status.Success {
			exitCode = status.EnableArrowKeyBindings
		}
	},
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

var disableArrowKeyBindingsCmd = &cobra.Command{
	Use:   "arrow_key_bindings",
	Short: "disable bindings for arrow keys (up/down)",
	Long: "Disable bindings for arrow keys (up/down)\n" +
		"Note that this only affects sessions of the same shell.\n" +
		"(e.g. running this in zsh will only enable the keybinding in zsh)",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = enableDisableArrowKeyBindingsGlobally(false)
		if exitCode == status.Success {
			exitCode = status.DisableArrowKeyBindings
		}
	},
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

func enableDisableArrowKeyBindingsGlobally(value bool) status.Code {
	usr, _ := user.Current()
	dir := usr.HomeDir
	configPath := filepath.Join(dir, ".config/resh.toml")
	var config cfg.Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		fmt.Println("Error reading config", err)
		return status.Fail
	}
	shell, found := os.LookupEnv("__RESH_ctl_shell")
	// shell env variable must be set and must be equal to either bash or zsh
	if found == false || (shell != "bash" && shell != "zsh") {
		fmt.Println("Error while determining a shell you are using - your RESH instalation is probably broken. Please reinstall RESH - exiting!")
		fmt.Println("found=", found, "shell=", shell)
		return status.Fail
	}
	if shell == "bash" {
		err := setConfigBindArrowKey(configPath, &config, &config.BindArrowKeysBash, shell, value)
		if err != nil {
			return status.Fail
		}
	} else if shell == "zsh" {
		err := setConfigBindArrowKey(configPath, &config, &config.BindArrowKeysZsh, shell, value)
		if err != nil {
			return status.Fail
		}
	} else {
		fmt.Println("FATAL ERROR while determining a shell you are using - your RESH instalation is probably broken. Please reinstall RESH - exiting!")
	}
	return status.Success
}

// I don't like the interface this function has - passing both config structure and a part of it feels wrong
// 		It's ugly and could lead to future errors
func setConfigBindArrowKey(configPath string, config *cfg.Config, configField *bool, shell string, value bool) error {
	if *configField != value {
		*configField = value

		f, err := os.Create(configPath)
		if err != nil {
			fmt.Println("Error: Failed to create/open file:", configPath, "; error:", err)
			return err
		}
		defer f.Close()
		if err := toml.NewEncoder(f).Encode(config); err != nil {
			fmt.Println("Error: Failed to encode and write the config values to hdd. error:", err)
			return err
		}
	}
	if value {
		fmt.Println("RESH arrow key bindings: ENABLED (in " + shell + ")")
	} else {
		fmt.Println("RESH arrow key bindings: DISABLED (in " + shell + ")")
	}
	return nil
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
