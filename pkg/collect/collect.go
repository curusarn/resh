package collect

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/curusarn/resh/pkg/httpclient"
	"github.com/curusarn/resh/pkg/records"
)

// SingleResponse json struct
type SingleResponse struct {
	Found   bool   `json:"found"`
	CmdLine string `json:"cmdline"`
}

// SendRecord to daemon
func SendRecord(r records.Record, port, path string) {
	recJSON, err := json.Marshal(r)
	if err != nil {
		log.Fatal("send err 1", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:"+port+path,
		bytes.NewBuffer(recJSON))
	if err != nil {
		log.Fatal("send err 2", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := httpclient.New()
	_, err = client.Do(req)
	if err != nil {
		log.Fatal("resh-daemon is not running - try restarting this terminal")
	}
}

// ReadFileContent and return it as a string
func ReadFileContent(path string) string {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
		//log.Fatal("failed to open " + path)
	}
	return strings.TrimSuffix(string(dat), "\n")
}

// GetGitDirs based on result of git "cdup" command
func GetGitDirs(cdup string, exitCode int, pwd string) (string, string) {
	if exitCode != 0 {
		return "", ""
	}
	abspath := filepath.Clean(filepath.Join(pwd, cdup))
	realpath, err := filepath.EvalSymlinks(abspath)
	if err != nil {
		log.Println("err while handling git dir paths:", err)
		return "", ""
	}
	return abspath, realpath
}

// GetTimezoneOffsetInSeconds based on zone returned by date command
func GetTimezoneOffsetInSeconds(zone string) float64 {
	// date +%z -> "+0200"
	hoursStr := zone[:3]
	minsStr := zone[3:]
	hours, err := strconv.Atoi(hoursStr)
	if err != nil {
		log.Println("err while parsing hours in timezone offset:", err)
		return -1
	}
	mins, err := strconv.Atoi(minsStr)
	if err != nil {
		log.Println("err while parsing mins in timezone offset:", err)
		return -1
	}
	secs := ((hours * 60) + mins) * 60
	return float64(secs)
}
