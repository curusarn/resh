package cmd

import (
	"github.com/curusarn/resh/cmd/control/status"
	"github.com/spf13/cobra"
)

var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "disable RESH features",
	Long:  `Disables RESH bindings for arrows and C-R.`,
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = status.DisableAll
	},
}

// var disableRecallingCmd = &cobra.Command{
// 	Use:   "keybind",
// 	Short: "Disables RESH bindings for arrows and C-R.",
// 	Long:  `Disables RESH bindings for arrows and C-R.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		exitCode = status.DisableAll
// 	},
// }
