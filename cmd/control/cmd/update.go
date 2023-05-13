package cmd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var betaFlag bool
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "check for updates and update RESH",
	Run: func(cmd *cobra.Command, args []string) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			out.FatalE("Could not get user home dir", err)
		}
		rawinstallPath := filepath.Join(homeDir, ".resh/rawinstall.sh")
		execArgs := []string{rawinstallPath}
		if betaFlag {
			execArgs = append(execArgs, "--beta")
		}
		execCmd := exec.Command("bash", execArgs...)
		execCmd.Stdin = os.Stdin
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		err = execCmd.Run()
		if err != nil {
			out.FatalE("Update ended with error", err)
		}
	},
}
