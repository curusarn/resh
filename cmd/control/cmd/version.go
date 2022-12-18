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
		printVersion("Installed", version, commit)

		versionEnv := getEnvVarWithDefault("__RESH_VERSION", "<unknown>")
		commitEnv := getEnvVarWithDefault("__RESH_REVISION", "<unknown>")
		printVersion("This terminal session", versionEnv, commitEnv)

		resp, err := status.GetDaemonStatus(config.Port)
		if err != nil {
			out.ErrorDaemonNotRunning(err)
			return
		}
		printVersion("Currently running daemon", resp.Version, resp.Commit)

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

func printVersion(title, version, commit string) {
	fmt.Printf("%s: %s (commit: %s)\n", title, version, commit)
}

func getEnvVarWithDefault(varName, defaultValue string) string {
	val, found := os.LookupEnv(varName)
	if !found {
		return defaultValue
	}
	return val
}
