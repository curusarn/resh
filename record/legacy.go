package record

type Legacy struct {
	// core
	CmdLine   string `json:"cmdLine"`
	ExitCode  int    `json:"exitCode"`
	Shell     string `json:"shell"`
	Uname     string `json:"uname"`
	SessionID string `json:"sessionId"`
	RecordID  string `json:"recordId"`

	// posix
	Home     string `json:"home"`
	Lang     string `json:"lang"`
	LcAll    string `json:"lcAll"`
	Login    string `json:"login"`
	Pwd      string `json:"pwd"`
	PwdAfter string `json:"pwdAfter"`
	ShellEnv string `json:"shellEnv"`
	Term     string `json:"term"`

	// non-posix"`
	RealPwd      string `json:"realPwd"`
	RealPwdAfter string `json:"realPwdAfter"`
	Pid          int    `json:"pid"`
	SessionPID   int    `json:"sessionPid"`
	Host         string `json:"host"`
	Hosttype     string `json:"hosttype"`
	Ostype       string `json:"ostype"`
	Machtype     string `json:"machtype"`
	Shlvl        int    `json:"shlvl"`

	// before after
	TimezoneBefore string `json:"timezoneBefore"`
	TimezoneAfter  string `json:"timezoneAfter"`

	RealtimeBefore      float64 `json:"realtimeBefore"`
	RealtimeAfter       float64 `json:"realtimeAfter"`
	RealtimeBeforeLocal float64 `json:"realtimeBeforeLocal"`
	RealtimeAfterLocal  float64 `json:"realtimeAfterLocal"`

	RealtimeDuration          float64 `json:"realtimeDuration"`
	RealtimeSinceSessionStart float64 `json:"realtimeSinceSessionStart"`
	RealtimeSinceBoot         float64 `json:"realtimeSinceBoot"`

	GitDir               string `json:"gitDir"`
	GitRealDir           string `json:"gitRealDir"`
	GitOriginRemote      string `json:"gitOriginRemote"`
	GitDirAfter          string `json:"gitDirAfter"`
	GitRealDirAfter      string `json:"gitRealDirAfter"`
	GitOriginRemoteAfter string `json:"gitOriginRemoteAfter"`
	MachineID            string `json:"machineId"`

	OsReleaseID         string `json:"osReleaseId"`
	OsReleaseVersionID  string `json:"osReleaseVersionId"`
	OsReleaseIDLike     string `json:"osReleaseIdLike"`
	OsReleaseName       string `json:"osReleaseName"`
	OsReleasePrettyName string `json:"osReleasePrettyName"`

	ReshUUID     string `json:"reshUuid"`
	ReshVersion  string `json:"reshVersion"`
	ReshRevision string `json:"reshRevision"`

	// records come in two parts (collect and postcollect)
	PartOne     bool `json:"partOne,omitempty"` // false => part two
	PartsMerged bool `json:"partsMerged"`
	// special flag -> not an actual record but an session end
	SessionExit bool `json:"sessionExit,omitempty"`

	// recall metadata
	Recalled          bool     `json:"recalled"`
	RecallHistno      int      `json:"recallHistno,omitempty"`
	RecallStrategy    string   `json:"recallStrategy,omitempty"`
	RecallActionsRaw  string   `json:"recallActionsRaw,omitempty"`
	RecallActions     []string `json:"recallActions,omitempty"`
	RecallLastCmdLine string   `json:"recallLastCmdLine"`

	// recall command
	RecallPrefix string `json:"recallPrefix,omitempty"`

	// added by sanitizatizer
	Sanitized bool `json:"sanitized,omitempty"`
	CmdLength int  `json:"cmdLength,omitempty"`

	// fields that are string here and int in older resh verisons
	Cols  interface{} `json:"cols"`
	Lines interface{} `json:"lines"`
}
