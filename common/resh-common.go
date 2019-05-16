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
    Pid int `json:"pid"`
    SessionPid int `json:"sessionPid"`
    WindowId int `json:"windowId"`
    Host string `json:"host"`
    Hosttype string `json:"hosttype"`
    Ostype string `json:"ostype"`
    Machtype string `json:"machtype"`

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

    GitWorkTree string `json:"gitWorkTree"`
}

type Config struct {
    Port int
}
