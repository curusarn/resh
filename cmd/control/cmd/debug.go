package cmd

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"

	"github.com/curusarn/resh/cmd/control/status"
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Shows logs and output from last runs of resh",
	Long:  "Shows logs and output from last runs of resh",
	Run: func(cmd *cobra.Command, args []string) {
		files := []string{
			"daemon_last_run_out.txt",
			"collect_last_run_out.txt",
			"postcollect_last_run_out.txt",
		}
		usr, _ := user.Current()
		dir := usr.HomeDir
		reshdir := filepath.Join(dir, ".resh")
		for _, fpath := range files {
			fpath := filepath.Join(reshdir, fpath)
			debugReadFile(fpath)
		}
		exitCode = status.Success
	},
}

func debugReadFile(path string) {
	fmt.Println("============================================================")
	fmt.Println(" filepath:", path)
	fmt.Println("============================================================")
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("ERROR while reading file:", err)
	}
	fmt.Println(string(dat))
}
