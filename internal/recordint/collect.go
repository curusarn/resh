package recordint

import "github.com/curusarn/resh/record"

type Collect struct {
	// record merging
	SessionID string
	Shlvl     int
	// session watching
	SessionPID int
	Shell      string

	Rec record.V1
}

type Postcollect struct {
	// record merging
	SessionID string
	Shlvl     int
	// session watching
	SessionPID int

	RecordID string
	ExitCode int
	Duration float64
}

type SessionInit struct {
	// record merging
	SessionID string
	Shlvl     int
	// session watching
	SessionPID int
}
