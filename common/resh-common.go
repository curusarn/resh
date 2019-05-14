package common

type Record struct {
    CmdLine string     `json: cmdLine`
    Pwd string         `json: pwd`
    GitWorkTree string `json: gitWorkTree`
    Shell string       `json: shell`
    ExitCode int       `json: exitCode`
    Logs []string      `json: logs`
}

