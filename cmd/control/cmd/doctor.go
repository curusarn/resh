package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/curusarn/resh/internal/cfg"
	"github.com/curusarn/resh/internal/check"
	"github.com/curusarn/resh/internal/msg"
	"github.com/curusarn/resh/internal/status"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func doctorCmdFunc(config cfg.Config) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		allOK := true
		if !checkDaemon(config) {
			allOK = false
			printDivider()
		}
		if !checkShellSession() {
			allOK = false
			printDivider()
		}
		if !checkShells() {
			allOK = false
			printDivider()
		}

		if allOK {
			out.Info("Everything looks good.")
		}
	}
}

func printDivider() {
	fmt.Printf("\n")
}

var msgFailedDaemonStart = `Failed to start RESH daemon.
 -> Start RESH daemon manually - run: resh-daemon-start
 -> Or restart this terminal window to bring RESH daemon back up
 -> You can check logs: ~/.local/share/resh/log.json (or ~/$XDG_DATA_HOME/resh/log.json)
 -> You can create an issue at: https://github.com/curusarn/resh/issues
`

func checkDaemon(config cfg.Config) bool {
	ok := true
	resp, err := status.GetDaemonStatus(config.Port)
	if err != nil {
		out.InfoE("RESH Daemon is not running", err)
		out.Info("Attempting to start RESH daemon ...")
		resp, err = startDaemon(config.Port, 5, 200*time.Millisecond)
		if err != nil {
			out.InfoE(msgFailedDaemonStart, err)
			return false
		}
		ok = false
		out.Info("Successfully started daemon.")
	}
	if version != resp.Version {
		out.InfoDaemonVersionMismatch(version, resp.Version)
		return false
	}
	return ok
}

func startDaemon(port int, maxRetries int, backoff time.Duration) (*msg.StatusResponse, error) {
	err := exec.Command("resh-daemon-start").Run()
	if err != nil {
		return nil, err
	}
	var resp *msg.StatusResponse
	retry := 0
	for {
		time.Sleep(backoff)
		resp, err = status.GetDaemonStatus(port)
		if err == nil {
			break
		}
		if retry == maxRetries {
			return nil, err
		}
		out.Logger.Error("Failed to get daemon status - retrying", zap.Error(err), zap.Int("retry", retry))
		retry++
		continue
	}
	return resp, nil
}

var msgShellFilesNotLoaded = `RESH shell files were not properly loaded in this terminal
 -> Try restarting this terminal to see if the issue persists
 -> Check your shell rc files (e.g. .zshrc, .bashrc, ...)
 -> You can create an issue at: https://github.com/curusarn/resh/issues
`

func checkShellSession() bool {
	versionEnv, found := os.LookupEnv("__RESH_VERSION")
	if !found {
		out.Info(msgShellFilesNotLoaded)
		return false
	}
	if version != versionEnv {
		out.InfoTerminalVersionMismatch(version, versionEnv)
		return false
	}
	return true
}

func checkShells() bool {
	allOK := true

	msg, err := check.LoginShell()
	if err != nil {
		out.InfoE("Failed to get login shell", err)
		allOK = false
	}
	if msg != "" {
		out.Info(msg)
		allOK = false
	}

	msg, err = check.ZshVersion()
	if err != nil {
		out.InfoE("Failed to check zsh version", err)
		allOK = false
	}
	if msg != "" {
		out.Info(msg)
		allOK = false
	}

	msg, err = check.BashVersion()
	if err != nil {
		out.InfoE("Failed to check bash version", err)
		allOK = false
	}
	if msg != "" {
		out.Info(msg)
		allOK = false
	}

	return allOK
}
