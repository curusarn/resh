package collect

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/curusarn/resh/internal/httpclient"
	"github.com/curusarn/resh/internal/output"
	"github.com/curusarn/resh/internal/records"
	"go.uber.org/zap"
)

// SingleResponse json struct
type SingleResponse struct {
	Found   bool   `json:"found"`
	CmdLine string `json:"cmdline"`
}

// SendRecord to daemon
func SendRecord(out *output.Output, r records.Record, port, path string) {
	out.Logger.Debug("Sending record ...",
		zap.String("cmdLine", r.CmdLine),
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

	client := httpclient.New()
	_, err = client.Do(req)
	if err != nil {
		out.FatalDaemonNotRunning(err)
	}
}

// ReadFileContent and return it as a string
func ReadFileContent(path string) string {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
		//sugar.Fatal("failed to open " + path)
	}
	return strings.TrimSuffix(string(dat), "\n")
}

// GetGitDirs based on result of git "cdup" command
func GetGitDirs(logger *zap.Logger, cdup string, exitCode int, pwd string) (string, string) {
	if exitCode != 0 {
		return "", ""
	}
	abspath := filepath.Clean(filepath.Join(pwd, cdup))
	realpath, err := filepath.EvalSymlinks(abspath)
	if err != nil {
		logger.Error("Error while handling git dir paths", zap.Error(err))
		return "", ""
	}
	return abspath, realpath
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
		logger.Error("err while parsing minutes in timezone offset:", zap.Error(err))
		return -1
	}
	secs := ((hours * 60) + mins) * 60
	return float64(secs)
}
