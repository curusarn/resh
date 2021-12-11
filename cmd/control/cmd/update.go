package cmd

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/curusarn/resh/cmd/control/status"
	"github.com/spf13/cobra"
)

var betaFlag bool
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "check for updates and update RESH",
	Run: func(cmd *cobra.Command, args []string) {
		usr, _ := user.Current()
		dir := usr.HomeDir
		rawinstallPath := filepath.Join(dir, ".resh/rawinstall.sh")
		execArgs := []string{rawinstallPath}
		if betaFlag {
			execArgs = append(execArgs, "--beta")
		}
		execCmd := exec.Command("bash", execArgs...)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		err := execCmd.Run()
		if err == nil {
			exitCode = status.Success
		}
	},
}
