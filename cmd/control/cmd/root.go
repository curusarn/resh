package cmd

import (
	"fmt"

	"github.com/curusarn/resh/internal/cfg"
	"github.com/curusarn/resh/internal/logger"
	"github.com/curusarn/resh/internal/output"
	"github.com/spf13/cobra"
)

var version string
var commit string

// globals
var out *output.Output

var rootCmd = &cobra.Command{
	Use:   "reshctl",
	Short: "Reshctl (RESH control) - check status, update",
}

// Execute reshctl
func Execute(ver, com, development string) {
	version = ver
	commit = com

	config, errCfg := cfg.New()
	logger, err := logger.New("reshctl", config.LogLevel, development)
	if err != nil {
		fmt.Printf("Error while creating logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any
	out = output.New(logger, "ERROR")
	if errCfg != nil {
		out.ErrorE("Error while getting configuration", errCfg)
	}

	var versionCmd = cobra.Command{
		Use:   "version",
		Short: "show RESH version",
		Run:   versionCmdFunc(config),
	}
	rootCmd.AddCommand(&versionCmd)

	doctorCmd := cobra.Command{
		Use:   "doctor",
		Short: "check common problems",
		Run:   doctorCmdFunc(config),
	}
	rootCmd.AddCommand(&doctorCmd)

	updateCmd.Flags().BoolVar(&betaFlag, "beta", false, "Update to latest version even if it's beta.")
	rootCmd.AddCommand(updateCmd)

	if err := rootCmd.Execute(); err != nil {
		out.FatalE("Command ended with error", err)
	}
}
