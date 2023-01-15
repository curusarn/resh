package cmd

import (
	"fmt"
	"os"

	"github.com/curusarn/resh/internal/cfg"
	"github.com/curusarn/resh/internal/status"
	"github.com/spf13/cobra"
)

func versionCmdFunc(config cfg.Config) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {

		fmt.Printf("Installed: %s\n", version)

		versionEnv := getEnvVarWithDefault("__RESH_VERSION", "<unknown>")
		fmt.Printf("This terminal session: %s\n", version)

		resp, err := status.GetDaemonStatus(config.Port)
		if err != nil {
			fmt.Printf("Running checks: %s\n", version)
			out.ErrorDaemonNotRunning(err)
			return
		}
		fmt.Printf("Currently running daemon: %s\n", resp.Version)

		if version != resp.Version {
			out.ErrorDaemonVersionMismatch(version, resp.Version)
			return
		}
		if version != versionEnv {
			out.ErrorTerminalVersionMismatch(version, versionEnv)
			return
		}
	}
}

func getEnvVarWithDefault(varName, defaultValue string) string {
	val, found := os.LookupEnv(varName)
	if !found {
		return defaultValue
	}
	return val
}
