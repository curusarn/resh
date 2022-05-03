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
	logger, _ := logger.New("postcollect", config.LogLevel, developement)
	defer logger.Sync() // flushes buffer, if any
	if errCfg != nil {
		logger.Error("Error while getting configuration", zap.Error(errCfg))
	}
	out := output.New(logger, "resh-postcollect ERROR")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		out.Fatal("Could not get user home dir", err)
	}
	reshUUIDPath := filepath.Join(homeDir, "/.resh/resh-uuid")
	machineIDPath := "/etc/machine-id"

	showVersion := flag.Bool("version", false, "Show version and exit")
	showRevision := flag.Bool("revision", false, "Show git revision and exit")

	requireVersion := flag.String("requireVersion", "", "abort if version doesn't match")
	requireRevision := flag.String("requireRevision", "", "abort if revision doesn't match")

	cmdLine := flag.String("cmdLine", "", "command line")
	exitCode := flag.Int("exitCode", -1, "exit code")
	sessionID := flag.String("sessionId", "", "resh generated session id")
	recordID := flag.String("recordId", "", "resh generated record id")

	shlvl := flag.Int("shlvl", -1, "$SHLVL")
	shell := flag.String("shell", "", "actual shell")

	// posix variables
	pwdAfter := flag.String("pwdAfter", "", "$PWD after command")

	// non-posix
	// sessionPid := flag.Int("sessionPid", -1, "$$ at session start")

	gitCdupAfter := flag.String("gitCdupAfter", "", "git rev-parse --show-cdup")
	gitRemoteAfter := flag.String("gitRemoteAfter", "", "git remote get-url origin")

	gitCdupExitCodeAfter := flag.Int("gitCdupExitCodeAfter", -1, "... $?")
	gitRemoteExitCodeAfter := flag.Int("gitRemoteExitCodeAfter", -1, "... $?")

	// before after
	timezoneAfter := flag.String("timezoneAfter", "", "")

	rtb := flag.String("realtimeBefore", "-1", "before $EPOCHREALTIME")
	rta := flag.String("realtimeAfter", "-1", "after $EPOCHREALTIME")
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
	realtimeAfter, err := strconv.ParseFloat(*rta, 64)
	if err != nil {
		out.Fatal("Error while parsing flag --realtimeAfter", err)
	}
	realtimeBefore, err := strconv.ParseFloat(*rtb, 64)
	if err != nil {
		out.Fatal("Error while parsing flag --realtimeBefore", err)
	}
	realtimeDuration := realtimeAfter - realtimeBefore

	timezoneAfterOffset := collect.GetTimezoneOffsetInSeconds(logger, *timezoneAfter)
	realtimeAfterLocal := realtimeAfter + timezoneAfterOffset

	realPwdAfter, err := filepath.EvalSymlinks(*pwdAfter)
	if err != nil {
		logger.Error("Error while handling pwdAfter realpath", zap.Error(err))
		realPwdAfter = ""
	}

	gitDirAfter, gitRealDirAfter := collect.GetGitDirs(logger, *gitCdupAfter, *gitCdupExitCodeAfter, *pwdAfter)
	if *gitRemoteExitCodeAfter != 0 {
		*gitRemoteAfter = ""
	}

	rec := records.Record{
		// core
		BaseRecord: records.BaseRecord{
			CmdLine:   *cmdLine,
			ExitCode:  *exitCode,
			SessionID: *sessionID,
			RecordID:  *recordID,
			Shlvl:     *shlvl,
			Shell:     *shell,

			PwdAfter: *pwdAfter,

			// non-posix
			RealPwdAfter: realPwdAfter,

			// before after
			TimezoneAfter: *timezoneAfter,

			RealtimeBefore:     realtimeBefore,
			RealtimeAfter:      realtimeAfter,
			RealtimeAfterLocal: realtimeAfterLocal,

			RealtimeDuration: realtimeDuration,

			GitDirAfter:          gitDirAfter,
			GitRealDirAfter:      gitRealDirAfter,
			GitOriginRemoteAfter: *gitRemoteAfter,
			MachineID:            collect.ReadFileContent(machineIDPath),

			PartOne: false,

			ReshUUID:     collect.ReadFileContent(reshUUIDPath),
			ReshVersion:  version,
			ReshRevision: commit,
		},
	}
	collect.SendRecord(out, rec, strconv.Itoa(config.Port), "/record")
}
