package cmd

import (
	"os"

	"github.com/curusarn/resh/cmd/control/status"
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash/zsh completion scripts",
	Long: `To load completion run

. <(reshctl completion bash) 

OR 

. <(reshctl completion zsh) 
`,
}

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

. <(reshctl completion bash) 
`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
		exitCode = status.Success
	},
}

var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generates zsh completion scripts",
	Long: `To load completion run

. <(reshctl completion zsh) 
`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenZshCompletion(os.Stdout)
		exitCode = status.Success
	},
}
