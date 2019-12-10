package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/curusarn/resh/pkg/cfg"
	"github.com/curusarn/resh/pkg/collect"
	"github.com/curusarn/resh/pkg/records"

	//  "os/exec"
	"os/user"
	"path/filepath"
	"strconv"
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

	var config cfg.Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Fatal("Error reading config:", err)
	}
	// recall command
	recall := flag.Bool("recall", false, "Recall command on position --histno")
	recallHistno := flag.Int("histno", 0, "Recall command on position --histno")
	recallPrefix := flag.String("prefix-search", "", "Recall command based on prefix --prefix-search")

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

	// recall metadata
	recallActions := flag.String("recall-actions", "", "recall actions that took place before executing the command")
	recallStrategy := flag.String("recall-strategy", "", "recall strategy used during recall actions")

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
	if *recallPrefix != "" && *recall == false {
		log.Println("Option '--prefix-search' only works with '--recall' option - exiting!")
		os.Exit(4)
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
	realtimeSinceSessionStart := realtimeBefore - realtimeSessionStart
	realtimeSinceBoot := realtimeSessSinceBoot + realtimeSinceSessionStart

	timezoneBeforeOffset := collect.GetTimezoneOffsetInSeconds(*timezoneBefore)
	realtimeBeforeLocal := realtimeBefore + timezoneBeforeOffset

	realPwd, err := filepath.EvalSymlinks(*pwd)
	if err != nil {
		log.Println("err while handling pwd realpath:", err)
		realPwd = ""
	}

	gitDir, gitRealDir := collect.GetGitDirs(*gitCdup, *gitCdupExitCode, *pwd)
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

	if *recall {
		rec := records.SlimRecord{
			SessionID:    *sessionID,
			RecallHistno: *recallHistno,
			RecallPrefix: *recallPrefix,
		}
		fmt.Print(collect.SendRecallRequest(rec, strconv.Itoa(config.Port)))
	} else {
		rec := records.Record{
			// posix
			Cols:  *cols,
			Lines: *lines,
			// core
			BaseRecord: records.BaseRecord{
				RecallHistno: *recallHistno,

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
				MachineID:       collect.ReadFileContent(machineIDPath),

				OsReleaseID:         *osReleaseID,
				OsReleaseVersionID:  *osReleaseVersionID,
				OsReleaseIDLike:     *osReleaseIDLike,
				OsReleaseName:       *osReleaseName,
				OsReleasePrettyName: *osReleasePrettyName,

				PartOne: true,

				ReshUUID:     collect.ReadFileContent(reshUUIDPath),
				ReshVersion:  Version,
				ReshRevision: Revision,

				RecallActionsRaw: *recallActions,
				RecallPrefix:     *recallPrefix,
				RecallStrategy:   *recallStrategy,
			},
		}
		collect.SendRecord(rec, strconv.Itoa(config.Port), "/record")
	}
}
