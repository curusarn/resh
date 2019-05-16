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
        log.Println("Error reading config", err)
        return
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

    host := flag.String("host", "", "$HOSTNAME")
    hosttype := flag.String("hosttype", "", "$HOSTTYPE")
    ostype := flag.String("ostype", "", "$OSTYPE")
    machtype := flag.String("machtype", "", "$MACHTYPE")

    // before after
    timezoneBefore := flag.String("timezoneBefore", "", "before $TZ")
    timezoneAfter := flag.String("timezoneAfter", "", "after $TZ")

    secondsUtcBefore := flag.Int("secsUtcBefore", -1, "secs utc")
    secondsUtcAfter := flag.Int("secsUtcAfter", -1, "secs utc")
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
    realtimeSessSinceBoot, err := strconv.ParseFloat(*rtsessboot, 64)
    realtimeDuration := realtimeAfter - realtimeBefore
    realtimeSinceSessionStart := realtimeBefore - realtimeSessionStart
    realtimeSinceBoot := realtimeSessSinceBoot + realtimeSinceSessionStart

    if err != nil {
        log.Fatal("Flag Parsing error:", err)
    }
    if err != nil {
        log.Fatal("Flag Parsing error:", err)
    }

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

        // before after
        TimezoneBefore: *timezoneBefore,
        TimezoneAfter: *timezoneAfter,

        RealtimeBefore: realtimeBefore,
        RealtimeAfter: realtimeAfter,
        SecondsUtcBefore: *secondsUtcBefore,
        SecondsUtcAfter: *secondsUtcAfter,
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
        log.Fatal("1 ", err)
    }

    req, err := http.NewRequest("POST", "http://localhost:" + port + "/record",
                                bytes.NewBuffer(recJson))
    if err != nil {
        log.Fatal("2 ", err)
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

