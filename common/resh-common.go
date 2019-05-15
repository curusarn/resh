package common

type Record struct {
    CmdLine string     `json:"cmdLine"`
    Pwd string         `json:"pwd"`
    GitWorkTree string `json:"gitWorkTree"`
    Shell string       `json:"shell"`
    ExitCode int       `json:"exitCode"`
    RealtimeBefore float64 `json:"realtimeBefore"`
    RealtimeAfter float64 `json:"realtimeAfter"`
    RealtimeDuration float64 `json:"realtimeDuration"`
    //Logs []string      `json: "logs"`
}

type Config struct {
    Port int
}
