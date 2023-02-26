package msg

import "github.com/curusarn/resh/internal/recordint"

// CliMsg struct
type CliMsg struct {
	SessionID string
	PWD       string
}

// CliResponse struct
type CliResponse struct {
	Records []recordint.SearchApp
}

// StatusResponse struct
type StatusResponse struct {
	Status  bool   `json:"status"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
}
