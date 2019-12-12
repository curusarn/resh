package cmd

import (
	"github.com/curusarn/resh/cmd/control/status"
	"github.com/spf13/cobra"
)

var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "disable RESH features (arrow key bindings)",
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
	Short: "disable bindings for arrow keys (up/down) FOR THIS AND ALL FUTURE SHELL SESSIONS",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: config set arrow_key_bindings true
		exitCode = status.DisableArrowKeyBindings
	},
}
