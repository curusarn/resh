package msg

// CliMsg struct
type CliMsg struct {
	SessionID string
	PWD       string
}

// CliResponse struct
type CliResponse struct {
	Records []record.SearchApp
}

// StatusResponse struct
type StatusResponse struct {
	Status  bool   `json:"status"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}
