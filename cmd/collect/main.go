package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/curusarn/resh/internal/cfg"
	"github.com/curusarn/resh/internal/collect"
	"github.com/curusarn/resh/internal/logger"
	"github.com/curusarn/resh/internal/output"
	"github.com/curusarn/resh/internal/records"
	"go.uber.org/zap"

	//  "os/exec"

	"path/filepath"
	"strconv"
)

// info passed during build
var version string
var commit string
var developement bool

func main() {
	config, errCfg := cfg.New()
	logger, _ := logger.New("collect", config.LogLevel, developement)
	defer logger.Sync() // flushes buffer, if any
	if errCfg != nil {
		logger.Error("Error while getting configuration", zap.Error(errCfg))
	}
	out := output.New(logger, "resh-collect ERROR")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		out.Fatal("Could not get user home dir", err)
	}

	reshUUIDPath := filepath.Join(homeDir, "/.resh/resh-uuid")
	machineIDPath := "/etc/machine-id"

	// version
	showVersion := flag.Bool("version", false, "Show version and exit")
	showRevision := flag.Bool("revision", false, "Show git revision and exit")

	requireVersion := flag.String("requireVersion", "", "abort if version doesn't match")
	requireRevision := flag.String("requireRevision", "", "abort if revision doesn't match")

	// core
	cmdLine := flag.String("cmdLine", "", "command line")
	exitCode := flag.Int("exitCode", -1, "exit code")
	shell := flag.String("shell", "", "actual shell")
	uname := flag.String("uname", "", "uname")
	sessionID := flag.String("sessionId", "", "resh generated session id")
	recordID := flag.String("recordId", "", "resh generated record id")

	// posix variables
	cols := flag.String("cols", "-1", "$COLUMNS")
	lines := flag.String("lines", "-1", "$LINES")
	home := flag.String("home", "", "$HOME")
	lang := flag.String("lang", "", "$LANG")
	lcAll := flag.String("lcAll", "", "$LC_ALL")
	login := flag.String("login", "", "$LOGIN")
	// path := flag.String("path", "", "$PATH")
	pwd := flag.String("pwd", "", "$PWD - present working directory")
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

	osReleaseID := flag.String("osReleaseId", "", "/etc/os-release ID")
	osReleaseVersionID := flag.String("osReleaseVersionId", "",
		"/etc/os-release ID")
	osReleaseIDLike := flag.String("osReleaseIdLike", "", "/etc/os-release ID")
	osReleaseName := flag.String("osReleaseName", "", "/etc/os-release ID")
	osReleasePrettyName := flag.String("osReleasePrettyName", "",
		"/etc/os-release ID")

	rtb := flag.String("realtimeBefore", "-1", "before $EPOCHREALTIME")
	rtsess := flag.String("realtimeSession", "-1",
		"on session start $EPOCHREALTIME")
	rtsessboot := flag.String("realtimeSessSinceBoot", "-1",
		"on session start $EPOCHREALTIME")
	flag.Parse()

	if *showVersion == true {
		fmt.Println(version)
		os.Exit(0)
	}
	if *showRevision == true {
		fmt.Println(commit)
		os.Exit(0)
	}
	if *requireVersion != "" && *requireVersion != version {
		fmt.Println("Please restart/reload this terminal session " +
			"(resh version: " + version +
			"; resh version of this terminal session: " + *requireVersion +
			")")
		os.Exit(3)
	}
	if *requireRevision != "" && *requireRevision != commit {
		fmt.Println("Please restart/reload this terminal session " +
			"(resh revision: " + commit +
			"; resh revision of this terminal session: " + *requireRevision +
			")")
		os.Exit(3)
	}

	realtimeBefore, err := strconv.ParseFloat(*rtb, 64)
	if err != nil {
		out.Fatal("Error while parsing flag --realtimeBefore", err)
	}
	realtimeSessionStart, err := strconv.ParseFloat(*rtsess, 64)
	if err != nil {
		out.Fatal("Error while parsing flag --realtimeSession", err)
	}
	realtimeSessSinceBoot, err := strconv.ParseFloat(*rtsessboot, 64)
	if err != nil {
		out.Fatal("Error while parsing flag --realtimeSessSinceBoot", err)
	}
	realtimeSinceSessionStart := realtimeBefore - realtimeSessionStart
	realtimeSinceBoot := realtimeSessSinceBoot + realtimeSinceSessionStart

	timezoneBeforeOffset := collect.GetTimezoneOffsetInSeconds(logger, *timezoneBefore)
	realtimeBeforeLocal := realtimeBefore + timezoneBeforeOffset

	realPwd, err := filepath.EvalSymlinks(*pwd)
	if err != nil {
		logger.Error("Error while handling pwd realpath", zap.Error(err))
		realPwd = ""
	}

	gitDir, gitRealDir := collect.GetGitDirs(logger, *gitCdup, *gitCdupExitCode, *pwd)
	if *gitRemoteExitCode != 0 {
		*gitRemote = ""
	}

	// if *osReleaseID == "" {
	// 	*osReleaseID = "linux"
	// }
	// if *osReleaseName == "" {
	// 	*osReleaseName = "Linux"
	// }
	// if *osReleasePrettyName == "" {
	// 	*osReleasePrettyName = "Linux"
	// }

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
			RecordID:  *recordID,

			// posix
			Home:  *home,
			Lang:  *lang,
			LcAll: *lcAll,
			Login: *login,
			// Path:     *path,
			Pwd:      *pwd,
			ShellEnv: *shellEnv,
			Term:     *term,

			// non-posix
			RealPwd:    realPwd,
			Pid:        *pid,
			SessionPID: *sessionPid,
			Host:       *host,
			Hosttype:   *hosttype,
			Ostype:     *ostype,
			Machtype:   *machtype,
			Shlvl:      *shlvl,

			// before after
			TimezoneBefore: *timezoneBefore,

			RealtimeBefore:      realtimeBefore,
			RealtimeBeforeLocal: realtimeBeforeLocal,

			RealtimeSinceSessionStart: realtimeSinceSessionStart,
			RealtimeSinceBoot:         realtimeSinceBoot,

			GitDir:          gitDir,
			GitRealDir:      gitRealDir,
			GitOriginRemote: *gitRemote,
			MachineID:       collect.ReadFileContent(out.Logger, machineIDPath),

			OsReleaseID:         *osReleaseID,
			OsReleaseVersionID:  *osReleaseVersionID,
			OsReleaseIDLike:     *osReleaseIDLike,
			OsReleaseName:       *osReleaseName,
			OsReleasePrettyName: *osReleasePrettyName,

			PartOne: true,

			ReshUUID:     collect.ReadFileContent(out.Logger, reshUUIDPath),
			ReshVersion:  version,
			ReshRevision: commit,
		},
	}
	collect.SendRecord(out, rec, strconv.Itoa(config.Port), "/record")
}
