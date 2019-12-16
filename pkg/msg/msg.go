package msg

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
