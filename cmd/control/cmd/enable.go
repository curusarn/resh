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

var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "enable RESH features (arrow key bindings)",
}

var enableArrowKeyBindingsCmd = &cobra.Command{
	Use:   "arrow_key_bindings",
	Short: "enable bindings for arrow keys (up/down) FOR THIS SHELL SESSION",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = status.EnableArrowKeyBindings
	},
}

var enableArrowKeyBindingsGlobalCmd = &cobra.Command{
	Use:   "arrow_key_bindings_global",
	Short: "enable bindings for arrow keys (up/down) FOR FUTURE SHELL SESSIONS",
	Long: "Enable bindings for arrow keys (up/down) FOR FUTURE SHELL SESSIONS.\n" +
		"Note that this only affects sessions of the same shell.\n" +
		"(e.g. running this in zsh will only affect future zsh sessions)",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = enableDisableArrowKeyBindingsGlobally(true)
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
	if *configField == value {
		if value {
			fmt.Println("The RESH arrow key bindings are ALREADY GLOBALLY ENABLED for all future " + shell + " sessions - nothing to do - exiting.")
		} else {
			fmt.Println("The RESH arrow key bindings are ALREADY GLOBALLY DISABLED for all future " + shell + " sessions - nothing to do - exiting.")
		}
		return nil
	}
	if value {
		fmt.Println("ENABLING the RESH arrow key bindings GLOBALLY (in " + shell + ") ...")
	} else {
		fmt.Println("DISABLING the RESH arrow key bindings GLOBALLY (in " + shell + ") ...")
	}
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
	if value {
		fmt.Println("SUCCESSFULLY ENABLED the RESH arrow key bindings GLOBALLY (in " + shell + ") " +
			"- every new (" + shell + ") session will start with enabled RESH arrow key bindings!")
	} else {
		fmt.Println("SUCCESSFULLY DISABLED the RESH arrow key bindings GLOBALLY (in " + shell + ") " +
			"- every new (" + shell + ") session will start with " + shell + " default arrow key bindings!")
	}
	return nil
}
