package common

type Record struct {
    // core
    CmdLine string     `json:"cmdLine"`
    ExitCode int       `json:"exitCode"`

    // posix
    Cols int `json:"cols"`
    Lines int `json:"lines"`
    Home string `json:"home"`
    Lang string `json:"lang"`
    LcAll string `json:"lcAll"`
    Login string `json:"login"`
    Path string `json:"path"`
    Pwd string `json:"pwd"`
    Shell string `json:"shell"`
    Term string `json:"term"`

    // non-posix"`
    RealPwd string `json:"realPwd"`
    Pid int `json:"pid"`
    ShellPid int `json:"shellPid"`
    WindowId int `json:"windowId"`
    Host string `json:"host"`
    Hosttype string `json:"hosttype"`
    Ostype string `json:"ostype"`
    Machtype string `json:"machtype"`
    Shlvl int `json:"shlvl"`

    // before after
    TimezoneBefore string `json:"timezoneBefore"`
    TimezoneAfter string `json:"timezoneAfter"`

    RealtimeBefore float64 `json:"realtimeBefore"`
    RealtimeAfter float64 `json:"realtimeAfter"`
    RealtimeBeforeLocal float64 `json:"realtimeBeforeLocal"`
    RealtimeAfterLocal float64 `json:"realtimeAfterLocal"`

    RealtimeDuration float64 `json:"realtimeDuration"`
    RealtimeSinceSessionStart float64 `json:"realtimeSinceSessionStart"`
    RealtimeSinceBoot float64 `json:"realtimeSinceBoot"`
    //Logs []string      `json: "logs"`

    GitDir string `json:"gitDir"`
    GitRealDir string `json:"gitRealDir"`
    GitOriginRemote string `json:"gitOriginRemote"`
    MachineId string `json:"machineId"`
    ReshUuid string `json:"reshUuid"`
}

type Config struct {
    Port int
}
