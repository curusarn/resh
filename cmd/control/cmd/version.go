package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/curusarn/resh/internal/msg"
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
			out.ErrorDaemonNotRunning(err)
			return
		}
		printVersion("Currently running daemon", resp.Version, resp.Commit)

		if version != resp.Version {
			out.ErrorDaemonVersionMismatch(version, resp.Version)
			return
		}
		if version != versionEnv {
			out.ErrorTerminalVersionMismatch(version, versionEnv)
			return
		}

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
	jsn, err := io.ReadAll(resp.Body)
	if err != nil {
		out.Fatal("Error while reading 'daemon /status' response", err)
	}
	err = json.Unmarshal(jsn, &mess)
	if err != nil {
		out.Fatal("Error while decoding 'daemon /status' response", err)
	}
	return mess, nil
}
