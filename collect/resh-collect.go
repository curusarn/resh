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

    cmdLine := flag.String("cmd", "", "command line")
    exitCode := flag.Int("status", -1, "exit code")
    pwd := flag.String("pwd", "", "present working directory")
    rtb := flag.String("realtimeBefore", "-1", "before $EPOCHREALTIME")
    rta := flag.String("realtimeAfter", "-1", "after $EPOCHREALTIME")
    flag.Parse()

    realtimeAfter, err := strconv.ParseFloat(*rta, 64)
    realtimeBefore, err := strconv.ParseFloat(*rtb, 64)
    realtimeDuration := realtimeAfter - realtimeBefore
    if err != nil {
        log.Fatal("Flag Parsing error:", err)
    }
    if err != nil {
        log.Fatal("Flag Parsing error:", err)
    }

    rec := common.Record{
        CmdLine: *cmdLine,
        Pwd: *pwd,
        GitWorkTree: getGitDir(),
        Shell: os.Getenv("SHELL"),
        ExitCode: *exitCode,
        RealtimeBefore: realtimeBefore,
        RealtimeAfter: realtimeAfter,
        RealtimeDuration: realtimeDuration,
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

