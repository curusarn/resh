package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/curusarn/resh/cmd/control/status"
	"github.com/curusarn/resh/pkg/msg"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show RESH version",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion("Installed", version, commit)

		versionEnv := getEnvVarWithDefault("__RESH_VERSION", "<unknown>")
		commitEnv := getEnvVarWithDefault("__RESH_REVISION", "<unknown>")
		printVersion("This terminal session", versionEnv, commitEnv)

		resp, err := getDaemonStatus(config.Port)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nERROR: Resh-daemon didn't respond - it's probably not running.\n\n")
			fmt.Fprintf(os.Stderr, "-> Try restarting this terminal window to bring resh-daemon back up.\n")
			fmt.Fprintf(os.Stderr, "-> If the problem persists you can check resh-daemon logs: ~/.resh/daemon.log\n")
			fmt.Fprintf(os.Stderr, "-> You can file an issue at: https://github.com/curusarn/resh/issues\n")
			exitCode = status.Fail
			return
		}
		printVersion("Currently running daemon", resp.Version, resp.Commit)

		if version != resp.Version {
			fmt.Fprintf(os.Stderr, "\nWARN: Resh-daemon is running in different version than is installed now - it looks like something went wrong during resh update.\n\n")
			fmt.Fprintf(os.Stderr, "-> Kill resh-daemon and then launch a new terminal window to fix that.\n")
			fmt.Fprintf(os.Stderr, " $ pkill resh-daemon\n")
			fmt.Fprintf(os.Stderr, "-> You can file an issue at: https://github.com/curusarn/resh/issues\n")
			return
		}
		if version != versionEnv {
			fmt.Fprintf(os.Stderr, "\nWARN: This terminal session was started with different resh version than is installed now - it looks like you updated resh and didn't restart this terminal.\n\n")
			fmt.Fprintf(os.Stderr, "-> Restart this terminal window to fix that.\n")
			return
		}

		exitCode = status.ReshStatus
	},
}

func printVersion(title, version, commit string) {
	fmt.Printf("%s: %s (commit: %s)\n", title, version, commit)
}

func getEnvVarWithDefault(varName, defaultValue string) string {
	val, found := os.LookupEnv(varName)
	if !found {
		return defaultValue
	}
	return val
}

func getDaemonStatus(port int) (msg.StatusResponse, error) {
	mess := msg.StatusResponse{}
	url := "http://localhost:" + strconv.Itoa(port) + "/status"
	resp, err := http.Get(url)
	if err != nil {
		return mess, err
	}
	defer resp.Body.Close()
	jsn, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error while reading 'daemon /status' response:", err)
	}
	err = json.Unmarshal(jsn, &mess)
	if err != nil {
		log.Fatal("Error while decoding 'daemon /status' response:", err)
	}
	return mess, nil
}
