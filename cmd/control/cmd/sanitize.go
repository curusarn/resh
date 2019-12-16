package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"

	"github.com/curusarn/resh/cmd/control/status"
	"github.com/spf13/cobra"
)

var sanitizeCmd = &cobra.Command{
	Use:   "sanitize",
	Short: "produce a sanitized version of your RESH history",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = status.Success
		usr, _ := user.Current()
		dir := usr.HomeDir

		fmt.Println()
		fmt.Println(" HOW IT WORKS")
		fmt.Println("   In sanitized history, all sensitive information is replaced with its SHA1 hashes.")
		fmt.Println()
		fmt.Println("Sanitizing ...")
		fmt.Println(" * ~/resh_history_sanitized.json (full lengh hashes)")
		execCmd := exec.Command("resh-sanitize", "-trim-hashes", "0", "--output", dir+"/resh_history_sanitized.json")
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		err := execCmd.Run()
		if err != nil {
			exitCode = status.Fail
		}

		fmt.Println(" * ~/resh_history_sanitized_trim12.json (12 char hashes)")
		execCmd = exec.Command("resh-sanitize", "-trim-hashes", "12", "--output", dir+"/resh_history_sanitized_trim12.json")
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		err = execCmd.Run()
		if err != nil {
			exitCode = status.Fail
		}
		fmt.Println()
		fmt.Println("Please direct all questions and/or issues to: https://github.com/curusarn/resh/issues")
		fmt.Println()
		fmt.Println("Please look at the resulting sanitized history using commands below.")
		fmt.Println(" * Pretty print JSON")
		fmt.Println("     cat ~/resh_history_sanitized_trim12.json | jq")
		fmt.Println()
		fmt.Println(" * Only show commands, don't show metadata")
		fmt.Println("     cat ~/resh_history_sanitized_trim12.json | jq '.[\"cmdLine\"]")
		fmt.Println()
	},
}
