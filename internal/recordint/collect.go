package recordint

import "github.com/curusarn/resh/internal/record"

type Collect struct {
	// record merging
	SessionID string
	Shlvl     int
	// session watching
	SessionPID int

	Rec record.V1
}

type Postcollect struct {
	// record merging
	SessionID string
	Shlvl     int
	// session watching
	SessionPID int
}
