package collect

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/curusarn/resh/internal/output"
	"github.com/curusarn/resh/internal/recordint"
	"go.uber.org/zap"
)

// SendRecord to daemon
func SendRecord(out *output.Output, r recordint.Collect, port, path string) {
	out.Logger.Debug("Sending record ...",
		zap.String("cmdLine", r.Rec.CmdLine),
		zap.String("sessionID", r.SessionID),
	)
	recJSON, err := json.Marshal(r)
	if err != nil {
		out.Fatal("Error while encoding record", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:"+port+path,
		bytes.NewBuffer(recJSON))
	if err != nil {
		out.Fatal("Error while sending record", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 1 * time.Second,
	}
	_, err = client.Do(req)
	if err != nil {
		out.FatalDaemonNotRunning(err)
	}
}

// SendSessionInit to daemon
func SendSessionInit(out *output.Output, r recordint.SessionInit, port string) {
	out.Logger.Debug("Sending session init ...",
		zap.String("sessionID", r.SessionID),
		zap.Int("sessionPID", r.SessionPID),
	)
	recJSON, err := json.Marshal(r)
	if err != nil {
		out.Fatal("Error while encoding record", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:"+port+"/session_init",
		bytes.NewBuffer(recJSON))
	if err != nil {
		out.Fatal("Error while sending record", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 1 * time.Second,
	}
	_, err = client.Do(req)
	if err != nil {
		out.FatalDaemonNotRunning(err)
	}
}

// ReadFileContent and return it as a string
func ReadFileContent(logger *zap.Logger, path string) string {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error("Error reading file",
			zap.String("filePath", path),
			zap.Error(err),
		)
		return ""
	}
	return strings.TrimSuffix(string(dat), "\n")
}

// GetGitDirs based on result of git "cdup" command
func GetGitDirs(logger *zap.Logger, cdUp string, exitCode int, pwd string) (string, string) {
	if exitCode != 0 {
		return "", ""
	}
	absPath := filepath.Clean(filepath.Join(pwd, cdUp))
	realPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		logger.Error("Error while handling git dir paths", zap.Error(err))
		return "", ""
	}
	return absPath, realPath
}

// GetTimezoneOffsetInSeconds based on zone returned by date command
func GetTimezoneOffsetInSeconds(logger *zap.Logger, zone string) float64 {
	// date +%z -> "+0200"
	hoursStr := zone[:3]
	minsStr := zone[3:]
	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		logger.Error("Error while parsing hours in timezone offset", zap.Error(err))
		return -1
	}
	mins, err := strconv.Atoi(minsStr)
	if err != nil {
		logger.Error("Errot while parsing minutes in timezone offset:", zap.Error(err))
		return -1
	}
	secs := ((hours * 60) + mins) * 60
	return float64(secs)
}
