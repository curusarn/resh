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
	showVersion := flag.Bool("version", false, "Show version and exit")
	showRevision := flag.Bool("revision", false, "Show git revision and exit")

	requireVersion := flag.String("requireVersion", "", "abort if version doesn't match")
	requireRevision := flag.String("requireRevision", "", "abort if revision doesn't match")

	cmdLine := flag.String("cmdLine", "", "command line")
	exitCode := flag.Int("exitCode", -1, "exit code")
	sessionID := flag.String("sessionId", "", "resh generated session id")

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
	realtimeDuration := realtimeAfter - realtimeBefore

	timezoneAfterOffset := collect.GetTimezoneOffsetInSeconds(*timezoneAfter)
	realtimeAfterLocal := realtimeAfter + timezoneAfterOffset

	realPwdAfter, err := filepath.EvalSymlinks(*pwdAfter)
	if err != nil {
		log.Println("err while handling pwdAfter realpath:", err)
		realPwdAfter = ""
	}

	gitDirAfter, gitRealDirAfter := collect.GetGitDirs(*gitCdupAfter, *gitCdupExitCodeAfter, *pwdAfter)
	if *gitRemoteExitCodeAfter != 0 {
		*gitRemoteAfter = ""
	}

	rec := records.Record{
		// core
		BaseRecord: records.BaseRecord{
			CmdLine:   *cmdLine,
			ExitCode:  *exitCode,
			SessionID: *sessionID,

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
			ReshVersion:  Version,
			ReshRevision: Revision,
		},
	}
	collect.SendRecord(rec, strconv.Itoa(config.Port), "/record")
}
