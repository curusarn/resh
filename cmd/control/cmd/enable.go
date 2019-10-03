package cmd

import (
	"github.com/curusarn/resh/cmd/control/status"
	"github.com/spf13/cobra"
)

var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "enable RESH features",
	Long:  `Enables RESH bindings for arrows and C-R.`,
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = status.EnableAll
	},
}

// var enableRecallingCmd = &cobra.Command{
// 	Use:   "keybind",
// 	Short: "Enables RESH bindings for arrows and C-R.",
// 	Long:  `Enables RESH bindings for arrows and C-R.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		exitCode = status.EnableAll
// 	},
// }
