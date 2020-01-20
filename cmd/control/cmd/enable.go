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

var enableControlRBindingCmd = &cobra.Command{
	Use:   "ctrl_r_binding",
	Short: "enable binding for control+R FOR THIS SHELL SESSION",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = status.EnableControlRBinding
	},
}

var enableControlRBindingGlobalCmd = &cobra.Command{
	Use:   "ctrl_r_binding_global",
	Short: "enable bindings for control+R FOR FUTURE SHELL SESSIONS",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = enableDisableArrowKeyBindingsGlobally(true)
	},
}

// Disable commands

var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "disable RESH features (bindings)",
}

var disableArrowKeyBindingsCmd = &cobra.Command{
	Use:   "arrow_key_bindings",
	Short: "disable bindings for arrow keys (up/down) FOR THIS SHELL SESSION",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = status.DisableArrowKeyBindings
	},
}

var disableArrowKeyBindingsGlobalCmd = &cobra.Command{
	Use:   "arrow_key_bindings_global",
	Short: "disable bindings for arrow keys (up/down) FOR FUTURE SHELL SESSIONS",
	Long: "Disable bindings for arrow keys (up/down) FOR FUTURE SHELL SESSIONS.\n" +
		"Note that this only affects sessions of the same shell.\n" +
		"(e.g. running this in zsh will only affect future zsh sessions)",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = enableDisableControlRBindingGlobally(false)
	},
}

var disableControlRBindingCmd = &cobra.Command{
	Use:   "ctrl_r_binding",
	Short: "disable binding for control+R FOR THIS SHELL SESSION",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = status.DisableControlRBinding
	},
}

var disableControlRBindingGlobalCmd = &cobra.Command{
	Use:   "ctrl_r_binding_global",
	Short: "disable bindings for control+R FOR FUTURE SHELL SESSIONS",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = enableDisableControlRBindingGlobally(false)
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

func enableDisableControlRBindingGlobally(value bool) status.Code {
	usr, _ := user.Current()
	dir := usr.HomeDir
	configPath := filepath.Join(dir, ".config/resh.toml")
	var config cfg.Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		fmt.Println("Error reading config", err)
		return status.Fail
	}
	if config.BindControlR == value {
		if value {
			fmt.Println("The RESH control+R binding is ALREADY GLOBALLY ENABLED for all future shell sessions - nothing to do - exiting.")
		} else {
			fmt.Println("The RESH control+R binding is ALREADY GLOBALLY DISABLED for all future shell sessions - nothing to do - exiting.")
		}
		return status.Fail
	}
	if value {
		fmt.Println("ENABLING the RESH arrow key bindings GLOBALLY ...")
	} else {
		fmt.Println("DISABLING the RESH arrow key bindings GLOBALLY ...")
	}
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
	if value {
		fmt.Println("SUCCESSFULLY ENABLED the RESH arrow key bindings GLOBALLY " +
			"- every new shell session will start with enabled RESH CLI control+R binding!")
	} else {
		fmt.Println("SUCCESSFULLY DISABLED the RESH arrow key bindings GLOBALLY " +
			"- every new shell session will start with your orignal control+R key binding!")
	}
	return status.Success
}
