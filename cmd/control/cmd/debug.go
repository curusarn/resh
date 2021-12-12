package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/curusarn/resh/cmd/control/status"
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "debug utils for resh",
	Long:  "Reloads resh rc files. Shows logs and output from last runs of resh",
}

var debugReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "reload resh rc files",
	Long:  "Reload resh rc files",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = status.ReloadRcFiles
	},
}

var debugInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "inspect session history",
	Run: func(cmd *cobra.Command, args []string) {
		exitCode = status.InspectSessionHistory
	},
}

var debugOutputCmd = &cobra.Command{
	Use:   "output",
	Short: "shows output from last runs of resh",
	Long:  "Shows output from last runs of resh",
	Run: func(cmd *cobra.Command, args []string) {
		files := []string{
			"daemon_last_run_out.txt",
			"collect_last_run_out.txt",
			"postcollect_last_run_out.txt",
			"session_init_last_run_out.txt",
			"cli_last_run_out.txt",
		}
		dir := os.Getenv("__RESH_XDG_CACHE_HOME")
		for _, fpath := range files {
			fpath := filepath.Join(dir, fpath)
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
