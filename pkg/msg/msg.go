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
