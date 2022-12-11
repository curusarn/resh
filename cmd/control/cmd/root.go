package cmd

import (
	"github.com/curusarn/resh/internal/cfg"
	"github.com/curusarn/resh/internal/logger"
	"github.com/curusarn/resh/internal/output"
	"github.com/spf13/cobra"
)

var version string
var commit string

// globals
var config cfg.Config
var out *output.Output

var rootCmd = &cobra.Command{
	Use:   "reshctl",
	Short: "Reshctl (RESH control) - check status, update, enable/disable features, sanitize history and more.",
}

// Execute reshctl
func Execute(ver, com, development string) {
	version = ver
	commit = com

	config, errCfg := cfg.New()
	logger, _ := logger.New("reshctl", config.LogLevel, development)
	defer logger.Sync() // flushes buffer, if any
	out = output.New(logger, "ERROR")
	if errCfg != nil {
		out.Error("Error while getting configuration", errCfg)
	}

	rootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(completionBashCmd)
	completionCmd.AddCommand(completionZshCmd)

	rootCmd.AddCommand(versionCmd)

	updateCmd.Flags().BoolVar(&betaFlag, "beta", false, "Update to latest version even if it's beta.")
	rootCmd.AddCommand(updateCmd)

	if err := rootCmd.Execute(); err != nil {
		out.Fatal("Command ended with error", err)
	}
}
