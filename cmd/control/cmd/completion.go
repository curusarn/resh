package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "generate bash/zsh completion scripts",
	Long: `To load completion run

. <(reshctl completion bash) 

OR 

. <(reshctl completion zsh) && compdef _reshctl reshctl
`,
}

var completionBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "generate bash completion scripts",
	Long: `To load completion run

. <(reshctl completion bash) 
`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
}

var completionZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "generate zsh completion scripts",
	Long: `To load completion run

. <(reshctl completion zsh) && compdef _reshctl reshctl
`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenZshCompletion(os.Stdout)
	},
}
