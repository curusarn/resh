package msg

import "github.com/curusarn/resh/pkg/records"

// CliMsg struct
type CliMsg struct {
	SessionID string `json:"sessionID"`
	PWD       string `json:"pwd"`
}

// CliResponse struct
type CliResponse struct {
	CliRecords []records.CliRecord `json:"cliRecords"`
}

// InspectMsg struct
type InspectMsg struct {
	SessionID string `json:"sessionId"`
	Count     uint   `json:"count"`
}

// MultiResponse struct
type MultiResponse struct {
	CmdLines []string `json:"cmdlines"`
}

// StatusResponse struct
type StatusResponse struct {
	Status  bool   `json:"status"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}
