package cmd

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"

	"github.com/curusarn/resh/cmd/control/status"
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug utils for resh",
	Long:  "Reloads resh rc files. Shows logs and output from last runs of resh",
}

var debugReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload resh rc files",
	Long:  "Reload resh rc files",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = status.ReloadRcFiles
	},
}

var debugOutputCmd = &cobra.Command{
	Use:   "output",
	Short: "Shows output from last runs of resh",
	Long:  "Shows output from last runs of resh",
	Run: func(cmd *cobra.Command, args []string) {
		files := []string{
			"daemon_last_run_out.txt",
			"collect_last_run_out.txt",
			"postcollect_last_run_out.txt",
			"session_init_last_run_out.txt",
			"arrow_up_last_run_out.txt",
			"arrow_down_last_run_out.txt",
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
