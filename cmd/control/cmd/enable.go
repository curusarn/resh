package cmd

import (
	"github.com/curusarn/resh/cmd/control/status"
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
	Short: "enable bindings for arrow keys (up/down) FOR THIS AND ALL FUTURE SHELL SESSIONS",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: config set arrow_key_bindings true
		exitCode = status.EnableArrowKeyBindings
	},
}
