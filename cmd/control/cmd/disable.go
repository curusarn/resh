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
	Short: "disable bindings for arrow keys (up/down) FOR FUTURE SHELL SESSIONS",
	Long: "Disable bindings for arrow keys (up/down) FOR FUTURE SHELL SESSIONS.\n" +
		"Note that this only affects sessions of the same shell.\n" +
		"(e.g. running this in zsh will only affect future zsh sessions)",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = enableDisableArrowKeyBindingsGlobally(false)
	},
}
