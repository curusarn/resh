package status

// Code - exit code of the resh-control command
type Code int

const (
	// Success exit code
	Success Code = 0
	// Fail exit code
	Fail = 1
	// EnableAll exit code - tells reshctl() wrapper to enable_all
	EnableAll = 100
	// DisableAll exit code - tells reshctl() wrapper to disable_all
	DisableAll = 110
	// ReloadRcFiles exit code - tells reshctl() wrapper to reload shellrc resh file
	ReloadRcFiles = 200
)
