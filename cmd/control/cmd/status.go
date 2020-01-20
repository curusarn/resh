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

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "show RESH status (aka systemctl status)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("resh " + version)
		fmt.Println()
		fmt.Println("Resh versions ...")
		fmt.Println(" * installed: " + version + " (" + commit + ")")
		resp, err := getDaemonStatus(config.Port)
		if err != nil {
			fmt.Println(" * daemon: NOT RUNNING!")
		} else {
			fmt.Println(" * daemon: " + resp.Version + " (" + resp.Commit + ")")
		}
		versionEnv, found := os.LookupEnv("__RESH_VERSION")
		if found == false {
			versionEnv = "UNKNOWN!"
		}
		commitEnv, found := os.LookupEnv("__RESH_REVISION")
		if found == false {
			commitEnv = "unknown"
		}
		fmt.Println(" * this session: " + versionEnv + " (" + commitEnv + ")")
		if version != resp.Version || version != versionEnv {
			fmt.Println(" * THERE IS A MISMATCH BETWEEN VERSIONS!")
			fmt.Println(" * Please REPORT this here: https://github.com/curusarn/resh/issues")
			fmt.Println(" * Please RESTART this terminal window")
		}

		fmt.Println()
		fmt.Println("Arrow key bindings ...")
		if config.BindArrowKeysBash {
			fmt.Println(" * bash future sessions: ENABLED (not recommended)")
		} else {
			fmt.Println(" * bash future sessions: DISABLED (recommended)")
		}
		if config.BindArrowKeysZsh {
			fmt.Println(" * zsh future sessions: ENABLED (recommended)")
		} else {
			fmt.Println(" * zsh future sessions: DISABLED (not recommended)")
		}

		exitCode = status.ReshStatus
	},
}

func getDaemonStatus(port int) (msg.StatusResponse, error) {
	mess := msg.StatusResponse{}
	url := "http://localhost:" + strconv.Itoa(port) + "/status"
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Daemon is not running!", err)
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
