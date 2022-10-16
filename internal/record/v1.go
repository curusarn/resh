package record

type V1 struct {
	// flags
	// deleted, favorite
	// FIXME: is this the best way? .. what about string, separate fields, or something similar
	Flags int `json:"flags"`

	DeviceID  string `json:"deviceID"`
	SessionID string `json:"sessionID"`
	// can we have a shorter uuid for record
	RecordID string `json:"recordID"`

	// cmdline, exitcode
	CmdLine  string `json:"cmdLine"`
	ExitCode int    `json:"exitCode"`

	// paths
	Home    string `json:"home"`
	Pwd     string `json:"pwd"`
	RealPwd string `json:"realPwd"`

	// hostname + lognem (not sure if we actually need logname)
	Logname  string `json:"logname"`
	Hostname string `json:"hostname"`

	// git info
	// origin is the most important
	GitOriginRemote string `json:"gitOriginRemote"`
	// maybe branch could be useful - e.g. in monorepo ??
	GitBranch string `json:"gitBranch"`

	// what is this for ??
	// session watching needs this
	// but I'm not sure if we need to save it
	// records belong to sessions
	// PID int `json:"pid"`
	// needed for tracking of sessions but I think it shouldn't be part of V1
	SessionPID int `json:"sessionPID"`

	// needed to because records are merged with parts with same "SessionID + Shlvl"
	// I don't think we need to save it
	Shlvl int `json:"shlvl"`

	// time (before), duration of command
	Time     float64 `json:"time"`
	Duration float64 `json:"duration"`

	// these look like internal stuff

	// records come in two parts (collect and postcollect)
	PartOne        bool `json:"partOne,omitempty"` // false => part two
	PartsNotMerged bool `json:"partsNotMerged,omitempty"`

	// special flag -> not an actual record but an session end
	// TODO: this shouldn't be part of serializable V1 record
	SessionExit bool `json:"sessionExit,omitempty"`
}
