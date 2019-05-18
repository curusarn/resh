package main

import (
    "bytes"
    "encoding/json"
    "os"
    "os/exec"
    "os/user"
    "path/filepath"
    "flag"
    "log"
    "net/http"
    "strconv"
    "strings"
    common "github.com/curusarn/resh/common"
    "github.com/BurntSushi/toml"
)

func main() {
    usr, _ := user.Current()
    dir := usr.HomeDir
    configPath := filepath.Join(dir, "/.config/resh.toml")

    var config common.Config
    if _, err := toml.DecodeFile(configPath, &config); err != nil {
        log.Fatal("Error reading config:", err)
    }

    cmdLine := flag.String("cmdLine", "", "command line")
    exitCode := flag.Int("exitCode", -1, "exit code")

    // posix variables
    cols := flag.Int("cols", -1, "$COLUMNS")
    lines := flag.Int("lines", -1, "$LINES")

    home := flag.String("home", "", "$HOME")
    lang := flag.String("lang", "", "$LANG")
    lcAll := flag.String("lcAll", "", "$LC_ALL")
    login := flag.String("login", "", "$LOGIN")
    path := flag.String("path", "", "$PATH")
    pwd := flag.String("pwd", "", "$PWD - present working directory")
    shell := flag.String("shell", "", "$SHELL")
    term := flag.String("term", "", "$TERM")

    // non-posix
    pid := flag.Int("pid", -1, "$PID")
    sessionPid := flag.Int("sessionPid", -1, "$$")
    windowId := flag.Int("windowId", -1, "$WINDOWID - session id")
    shlvl := flag.Int("shlvl", -1, "$SHLVL")

    host := flag.String("host", "", "$HOSTNAME")
    hosttype := flag.String("hosttype", "", "$HOSTTYPE")
    ostype := flag.String("ostype", "", "$OSTYPE")
    machtype := flag.String("machtype", "", "$MACHTYPE")

    // before after
    timezoneBefore := flag.String("timezoneBefore", "", "")
    timezoneAfter := flag.String("timezoneAfter", "", "")

    rtb := flag.String("realtimeBefore", "-1", "before $EPOCHREALTIME")
    rta := flag.String("realtimeAfter", "-1", "after $EPOCHREALTIME")
    rtsess := flag.String("realtimeSession", "-1",
                          "on session start $EPOCHREALTIME")
    rtsessboot := flag.String("realtimeSessSinceBoot", "-1",
                              "on session start $EPOCHREALTIME")
    flag.Parse()

    realtimeAfter, err := strconv.ParseFloat(*rta, 64)
    realtimeBefore, err := strconv.ParseFloat(*rtb, 64)
    realtimeSessionStart, err := strconv.ParseFloat(*rtsess, 64)
    if err != nil {
        log.Fatal("Flag Parsing error (1):", err)
    }
    realtimeSessSinceBoot, err := strconv.ParseFloat(*rtsessboot, 64)
    if err != nil {
        log.Fatal("Flag Parsing error (2):", err)
    }
    realtimeDuration := realtimeAfter - realtimeBefore
    realtimeSinceSessionStart := realtimeBefore - realtimeSessionStart
    realtimeSinceBoot := realtimeSessSinceBoot + realtimeSinceSessionStart

    timezoneBeforeOffset := getTimezoneOffsetInSeconds(*timezoneBefore)
    timezoneAfterOffset := getTimezoneOffsetInSeconds(*timezoneAfter)
    realtimeBeforeLocal := realtimeBefore + timezoneBeforeOffset
    realtimeAfterLocal := realtimeAfter + timezoneAfterOffset


    rec := common.Record{
        // core
        CmdLine: *cmdLine,
        ExitCode: *exitCode,

        // posix
        Cols: *cols,
        Lines: *lines,

        Home: *home,
        Lang: *lang,
        LcAll: *lcAll,
        Login: *login,
        Path: *path,
        Pwd: *pwd,
        Shell: *shell,
        Term: *term,

        // non-posix
        Pid: *pid,
        SessionPid: *sessionPid,
        WindowId: *windowId,
        Host: *host,
        Hosttype: *hosttype,
        Ostype: *ostype,
        Machtype: *machtype,
        Shlvl: *shlvl,

        // before after
        TimezoneBefore: *timezoneBefore,
        TimezoneAfter: *timezoneAfter,

        RealtimeBefore: realtimeBefore,
        RealtimeAfter: realtimeAfter,
        RealtimeBeforeLocal: realtimeBeforeLocal,
        RealtimeAfterLocal: realtimeAfterLocal,

        RealtimeDuration: realtimeDuration,
        RealtimeSinceSessionStart: realtimeSinceSessionStart,
        RealtimeSinceBoot: realtimeSinceBoot,

        GitWorkTree: getGitDir(),
    }
    sendRecord(rec, strconv.Itoa(config.Port))
}

func sendRecord(r common.Record, port string) {
    recJson, err := json.Marshal(r)
    if err != nil {
        log.Fatal("send err 1", err)
    }

    req, err := http.NewRequest("POST", "http://localhost:" + port + "/record",
                                bytes.NewBuffer(recJson))
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

func getGitDir() string {
    // assume we are in pwd
    gitWorkTree := os.Getenv("GIT_WORK_TREE")

    if gitWorkTree != "" {
        return gitWorkTree
    }
    // we should look up the git directory ourselves
    // OR leave it to resh daemon to not slow down user
    out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
    if err != nil {
        if exitError, ok := err.(*exec.ExitError); ok {
            if exitError.ExitCode() == 128 {
                return ""
            }
            log.Fatal("git cmd failed")
        } else {
            log.Fatal("git cmd failed w/o exit code")
        }
    }
    return strings.TrimSuffix(string(out), "\n")
}

func getTimezoneOffsetInSeconds(zone string) float64 {
    hours_mins := strings.Split(zone, ":")
    hours, err := strconv.Atoi(hours_mins[0])
    if err != nil {
        log.Println("err while parsing hours in timezone offset:", err)
        return -1
    }
    mins, err := strconv.Atoi(hours_mins[1])
    if err != nil {
        log.Println("err while parsing mins in timezone offset:", err)
        return -1
    }
    secs := ( (hours * 60) + mins ) * 60
    return float64(secs)
}

