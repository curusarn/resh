package cmd

import (
	"os"
	"os/exec"

	"github.com/curusarn/resh/cmd/control/status"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "checks for updates and updates RESH",
	Run: func(cmd *cobra.Command, args []string) {
		url := "https://raw.githubusercontent.com/curusarn/resh/master/scripts/rawinstall.sh"
		execCmd := exec.Command("bash", "-c", "curl -fsSL "+url+" | bash")
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		err := execCmd.Run()
		if err == nil {
			exitCode = status.Success
		}
	},
}
