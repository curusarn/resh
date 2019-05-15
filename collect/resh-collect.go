package main

import (
    "bytes"
    "encoding/json"
    "os"
    "os/exec"
    "os/user"
    "path/filepath"
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

    exitCode, err := strconv.Atoi(os.Args[1])
    if err != nil {
        // log this 
        log.Fatal("First arg is not a number! (expecting $?)", err)
    }
    pwd := os.Args[2]
    cmdLine := os.Args[3]
    rec := common.Record{
        CmdLine: cmdLine,
        Pwd: pwd,
        GitWorkTree: getGitDir(),
        Shell: os.Getenv("SHELL"),
        ExitCode: exitCode,
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

