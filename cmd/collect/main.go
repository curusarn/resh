package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/curusarn/resh/pkg/records"

	//  "os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// Version from git set during build
var Version string

// Revision from git set during build
var Revision string

func main() {
	usr, _ := user.Current()
	dir := usr.HomeDir
	configPath := filepath.Join(dir, "/.config/resh.toml")
	reshUUIDPath := filepath.Join(dir, "/.resh/resh-uuid")

	machineIDPath := "/etc/machine-id"

	var config records.Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Fatal("Error reading config:", err)
	}
	showVersion := flag.Bool("version", false, "Show version and exit")
	showRevision := flag.Bool("revision", false, "Show git revision and exit")

	requireVersion := flag.String("requireVersion", "", "abort if version doesn't match")
	requireRevision := flag.String("requireRevision", "", "abort if revision doesn't match")

	cmdLine := flag.String("cmdLine", "", "command line")
	exitCode := flag.Int("exitCode", -1, "exit code")
	shell := flag.String("shell", "", "actual shell")
	uname := flag.String("uname", "", "uname")
	sessionID := flag.String("sessionId", "", "resh generated session id")

	// posix variables
	cols := flag.String("cols", "-1", "$COLUMNS")
	lines := flag.String("lines", "-1", "$LINES")
	home := flag.String("home", "", "$HOME")
	lang := flag.String("lang", "", "$LANG")
	lcAll := flag.String("lcAll", "", "$LC_ALL")
	login := flag.String("login", "", "$LOGIN")
	// path := flag.String("path", "", "$PATH")
	pwd := flag.String("pwd", "", "$PWD - present working directory")
	pwdAfter := flag.String("pwdAfter", "", "$PWD after command")
	shellEnv := flag.String("shellEnv", "", "$SHELL")
	term := flag.String("term", "", "$TERM")

	// non-posix
	pid := flag.Int("pid", -1, "$$")
	sessionPid := flag.Int("sessionPid", -1, "$$ at session start")
	shlvl := flag.Int("shlvl", -1, "$SHLVL")

	host := flag.String("host", "", "$HOSTNAME")
	hosttype := flag.String("hosttype", "", "$HOSTTYPE")
	ostype := flag.String("ostype", "", "$OSTYPE")
	machtype := flag.String("machtype", "", "$MACHTYPE")
	gitCdup := flag.String("gitCdup", "", "git rev-parse --show-cdup")
	gitRemote := flag.String("gitRemote", "", "git remote get-url origin")

	gitCdupExitCode := flag.Int("gitCdupExitCode", -1, "... $?")
	gitRemoteExitCode := flag.Int("gitRemoteExitCode", -1, "... $?")

	// before after
	timezoneBefore := flag.String("timezoneBefore", "", "")
	timezoneAfter := flag.String("timezoneAfter", "", "")

	osReleaseID := flag.String("osReleaseId", "", "/etc/os-release ID")
	osReleaseVersionID := flag.String("osReleaseVersionId", "",
		"/etc/os-release ID")
	osReleaseIDLike := flag.String("osReleaseIdLike", "", "/etc/os-release ID")
	osReleaseName := flag.String("osReleaseName", "", "/etc/os-release ID")
	osReleasePrettyName := flag.String("osReleasePrettyName", "",
		"/etc/os-release ID")

	rtb := flag.String("realtimeBefore", "-1", "before $EPOCHREALTIME")
	rta := flag.String("realtimeAfter", "-1", "after $EPOCHREALTIME")
	rtsess := flag.String("realtimeSession", "-1",
		"on session start $EPOCHREALTIME")
	rtsessboot := flag.String("realtimeSessSinceBoot", "-1",
		"on session start $EPOCHREALTIME")
	flag.Parse()

	if *showVersion == true {
		fmt.Println(Version)
		os.Exit(0)
	}
	if *showRevision == true {
		fmt.Println(Revision)
		os.Exit(0)
	}
	if *requireVersion != "" && *requireVersion != Version {
		fmt.Println("Please restart/reload this terminal session " +
			"(resh version: " + Version +
			"; resh version of this terminal session: " + *requireVersion +
			")")
		os.Exit(3)
	}
	if *requireRevision != "" && *requireRevision != Revision {
		fmt.Println("Please restart/reload this terminal session " +
			"(resh revision: " + Revision +
			"; resh revision of this terminal session: " + *requireRevision +
			")")
		os.Exit(3)
	}
	realtimeAfter, err := strconv.ParseFloat(*rta, 64)
	if err != nil {
		log.Fatal("Flag Parsing error (rta):", err)
	}
	realtimeBefore, err := strconv.ParseFloat(*rtb, 64)
	if err != nil {
		log.Fatal("Flag Parsing error (rtb):", err)
	}
	realtimeSessionStart, err := strconv.ParseFloat(*rtsess, 64)
	if err != nil {
		log.Fatal("Flag Parsing error (rt sess):", err)
	}
	realtimeSessSinceBoot, err := strconv.ParseFloat(*rtsessboot, 64)
	if err != nil {
		log.Fatal("Flag Parsing error (rt sess boot):", err)
	}
	realtimeDuration := realtimeAfter - realtimeBefore
	realtimeSinceSessionStart := realtimeBefore - realtimeSessionStart
	realtimeSinceBoot := realtimeSessSinceBoot + realtimeSinceSessionStart

	timezoneBeforeOffset := getTimezoneOffsetInSeconds(*timezoneBefore)
	timezoneAfterOffset := getTimezoneOffsetInSeconds(*timezoneAfter)
	realtimeBeforeLocal := realtimeBefore + timezoneBeforeOffset
	realtimeAfterLocal := realtimeAfter + timezoneAfterOffset

	realPwd, err := filepath.EvalSymlinks(*pwd)
	if err != nil {
		log.Println("err while handling pwd realpath:", err)
		realPwd = ""
	}
	realPwdAfter, err := filepath.EvalSymlinks(*pwdAfter)
	if err != nil {
		log.Println("err while handling pwdAfter realpath:", err)
		realPwd = ""
	}

	gitDir, gitRealDir := getGitDirs(*gitCdup, *gitCdupExitCode, *pwd)
	if *gitRemoteExitCode != 0 {
		*gitRemote = ""
	}

	if *osReleaseID == "" {
		*osReleaseID = "linux"
	}
	if *osReleaseName == "" {
		*osReleaseName = "Linux"
	}
	if *osReleasePrettyName == "" {
		*osReleasePrettyName = "Linux"
	}

	rec := records.Record{
		// posix
		Cols:  *cols,
		Lines: *lines,
		// core
		BaseRecord: records.BaseRecord{
			CmdLine:   *cmdLine,
			ExitCode:  *exitCode,
			Shell:     *shell,
			Uname:     *uname,
			SessionID: *sessionID,

			// posix
			Home:  *home,
			Lang:  *lang,
			LcAll: *lcAll,
			Login: *login,
			// Path:     *path,
			Pwd:      *pwd,
			PwdAfter: *pwdAfter,
			ShellEnv: *shellEnv,
			Term:     *term,

			// non-posix
			RealPwd:      realPwd,
			RealPwdAfter: realPwdAfter,
			Pid:          *pid,
			SessionPid:   *sessionPid,
			Host:         *host,
			Hosttype:     *hosttype,
			Ostype:       *ostype,
			Machtype:     *machtype,
			Shlvl:        *shlvl,

			// before after
			TimezoneBefore: *timezoneBefore,
			TimezoneAfter:  *timezoneAfter,

			RealtimeBefore:      realtimeBefore,
			RealtimeAfter:       realtimeAfter,
			RealtimeBeforeLocal: realtimeBeforeLocal,
			RealtimeAfterLocal:  realtimeAfterLocal,

			RealtimeDuration:          realtimeDuration,
			RealtimeSinceSessionStart: realtimeSinceSessionStart,
			RealtimeSinceBoot:         realtimeSinceBoot,

			GitDir:          gitDir,
			GitRealDir:      gitRealDir,
			GitOriginRemote: *gitRemote,
			MachineID:       readFileContent(machineIDPath),

			OsReleaseID:         *osReleaseID,
			OsReleaseVersionID:  *osReleaseVersionID,
			OsReleaseIDLike:     *osReleaseIDLike,
			OsReleaseName:       *osReleaseName,
			OsReleasePrettyName: *osReleasePrettyName,

			ReshUUID:     readFileContent(reshUUIDPath),
			ReshVersion:  Version,
			ReshRevision: Revision,
		},
	}
	sendRecord(rec, strconv.Itoa(config.Port))
}

func sendRecord(r records.Record, port string) {
	recJSON, err := json.Marshal(r)
	if err != nil {
		log.Fatal("send err 1", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:"+port+"/record",
		bytes.NewBuffer(recJSON))
	if err != nil {
		log.Fatal("send err 2", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Fatal("resh-daemon is not running :(")
	}
}

func readFileContent(path string) string {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
		//log.Fatal("failed to open " + path)
	}
	return strings.TrimSuffix(string(dat), "\n")
}

func getGitDirs(cdup string, exitCode int, pwd string) (string, string) {
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

func getTimezoneOffsetInSeconds(zone string) float64 {
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

// func getGitRemote() string {
// 	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
// 	if err != nil {
// 		if exitError, ok := err.(*exec.ExitError); ok {
// 			if exitError.ExitCode() == 128 {
// 				return ""
// 			}
// 			log.Fatal("git remote cmd failed")
// 		} else {
// 			log.Fatal("git remote cmd failed w/o exit code")
// 		}
// 	}
// 	return strings.TrimSuffix(string(out), "\n")
// }
//
// func getGitDir() string {
// 	// assume we are in pwd
// 	gitWorkTree := os.Getenv("GIT_WORK_TREE")
//
// 	if gitWorkTree != "" {
// 		return gitWorkTree
// 	}
// 	// we should look up the git directory ourselves
// 	// OR leave it to resh daemon to not slow down user
// 	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
// 	if err != nil {
// 		if exitError, ok := err.(*exec.ExitError); ok {
// 			if exitError.ExitCode() == 128 {
// 				return ""
// 			}
// 			log.Fatal("git rev-parse cmd failed")
// 		} else {
// 			log.Fatal("git rev-parse cmd failed w/o exit code")
// 		}
// 	}
// 	return strings.TrimSuffix(string(out), "\n")
// }
// }
