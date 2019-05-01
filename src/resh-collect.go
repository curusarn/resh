package main

import (
    "fmt"
    "os"
    "os/exec"
    "log"
    "strconv"
    "strings"
)

type record struct {
    CmdLine string
    Pwd string
    GitWorkTree string
    Shell string
    ExitCode int
    //Logs string[]
}

func main() {
    exitCode, err := strconv.Atoi(os.Args[1])
    if err != nil {
        // log this 
        log.Fatal("First arg is not a number! (expecting $?)", err)
    }
    pwd := os.Args[2]
    cmdLine := os.Args[3]
    rec := record{
        CmdLine: cmdLine,
        Pwd: pwd,
        GitWorkTree: getGitDir(),
        Shell: os.Getenv("SHELL"),
        ExitCode: exitCode,
    }
    rec.send()
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

func (r record) send() {
    fmt.Println("cmd:", r.CmdLine)
    fmt.Println("pwd:", r.Pwd)
    fmt.Println("git:", r.GitWorkTree)
    fmt.Println("exit_code:", r.ExitCode)
}
