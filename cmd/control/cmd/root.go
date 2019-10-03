package cmd

import (
	"fmt"
	"log"

	"github.com/curusarn/resh/cmd/control/status"
	"github.com/spf13/cobra"
)

var exitCode status.Code = status.DefaultInvalid

var rootCmd = &cobra.Command{
	Use:   "reshctl",
	Short: "Reshctl (RESH control) - enables you to enable/disable features and more.",
	Long:  `Enables you to enable/disable RESH bindings for arrows and C-R.`,
}

// Execute reshctl
func Execute() status.Code {
	rootCmd.AddCommand(disableCmd)
	// disableCmd.AddCommand(disableRecallingCmd)

	rootCmd.AddCommand(enableCmd)
	// enableCmd.AddCommand(enableRecallingCmd)

	rootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(completionBashCmd)
	completionCmd.AddCommand(completionZshCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return status.Fail
	}
	if exitCode == status.DefaultInvalid {
		log.Println("reshctl FATAL ERROR: (sub)command didn't set exitCode!")
		return status.Fail
	}
	return exitCode
}
