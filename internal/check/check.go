package check

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func LoginShell() (string, error) {
	shellPath, found := os.LookupEnv("SHELL")
	if !found {
		return "", fmt.Errorf("env variable $SHELL is not set")
	}
	parts := strings.Split(shellPath, "/")
	shell := parts[len(parts)-1]
	if shell != "bash" && shell != "zsh" {
		return fmt.Sprintf("Current shell (%s) is unsupported\n", shell), nil
	}
	return "", nil
}

func msgShellVersion(shell, expectedVer, actualVer string) string {
	return fmt.Sprintf(
		"Minimal supported %s version is %s. You have %s.\n"+
			" -> Update to %s %s+ if you want to use RESH with it",
		shell, expectedVer, actualVer,
		shell, expectedVer,
	)
}

func BashVersion() (string, error) {
	out, err := exec.Command("bash", "-c", "echo $BASH_VERSION").Output()
	if err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}
	verStr := strings.TrimSuffix(string(out), "\n")
	ver, err := parseVersion(verStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse version: %w", err)
	}

	if ver.Major < 4 || (ver.Major == 4 && ver.Minor < 3) {
		return msgShellVersion("bash", "4.3", verStr), nil
	}
	return "", nil
}

func ZshVersion() (string, error) {
	out, err := exec.Command("zsh", "-c", "echo $ZSH_VERSION").Output()
	if err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}
	verStr := strings.TrimSuffix(string(out), "\n")
	ver, err := parseVersion(string(out))
	if err != nil {
		return "", fmt.Errorf("failed to parse version: %w", err)
	}

	if ver.Major < 5 {
		return msgShellVersion("zsh", "5.0", verStr), nil
	}
	return "", nil
}

type version struct {
	Major int
	Minor int
	Rest  string
}

func parseVersion(str string) (version, error) {
	parts := strings.SplitN(str, ".", 3)
	if len(parts) < 3 {
		return version{}, fmt.Errorf("not enough parts")
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return version{}, fmt.Errorf("failed to parse major version: %w", err)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return version{}, fmt.Errorf("failed to parse minor version: %w", err)
	}
	ver := version{
		Major: major,
		Minor: minor,
		Rest:  parts[2],
	}
	return ver, nil
}
