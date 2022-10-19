package record

type V1 struct {
	// flags
	// deleted, favorite
	// FIXME: is this the best way? .. what about string, separate fields, or something similar
	Flags int `json:"flags"`

	// cmdline, exitcode
	CmdLine  string `json:"cmdLine"`
	ExitCode int    `json:"exitCode"`

	DeviceID  string `json:"deviceID"`
	SessionID string `json:"sessionID"`
	// can we have a shorter uuid for record
	RecordID string `json:"recordID"`

	// paths
	// TODO: Do we need both pwd and real pwd?
	Home    string `json:"home"`
	Pwd     string `json:"pwd"`
	RealPwd string `json:"realPwd"`

	// hostname + logname (not sure if we actually need logname)
	// Logname  string `json:"logname"`
	// Device is usually hostname but not stricly hostname
	// It can be configured in RESH configuration
	Device string `json:"device"`

	// git info
	// origin is the most important
	GitOriginRemote string `json:"gitOriginRemote"`
	// TODO: add GitBranch (v2 ?)
	//       maybe branch could be useful - e.g. in monorepo ??
	// GitBranch string `json:"gitBranch"`

	// what is this for ??
	// session watching needs this
	// but I'm not sure if we need to save it
	// records belong to sessions
	// PID int `json:"pid"`
	// needed for tracking of sessions but I think it shouldn't be part of V1
	// SessionPID int `json:"sessionPID"`

	// needed to because records are merged with parts with same "SessionID + Shlvl"
	// I don't think we need to save it
	// Shlvl int `json:"shlvl"`

	// time (before), duration of command
	// time and duration are strings because we don't want unnecessary precision when they get serialized into json
	// we could implement custom (un)marshalling but I don't see downsides of directly representing the values as strings
	Time     string `json:"time"`
	Duration string `json:"duration"`

	// these look like internal stuff

	// records come in two parts (collect and postcollect)
	PartOne        bool `json:"partOne,omitempty"` // false => part two
	PartsNotMerged bool `json:"partsNotMerged,omitempty"`

	// special flag -> not an actual record but an session end
	// TODO: this shouldn't be part of serializable V1 record
	SessionExit bool `json:"sessionExit,omitempty"`
}
