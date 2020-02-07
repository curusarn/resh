package status

// Code - exit code of the resh-control command
type Code int

const (
	// Success exit code
	Success Code = 0
	// Fail exit code
	Fail = 1
	// EnableResh exit code -  tells reshctl() wrapper to enable resh
	// EnableResh = 30

	// EnableArrowKeyBindings exit code - tells reshctl() wrapper to enable arrow key bindings
	EnableArrowKeyBindings = 31
	// EnableControlRBinding exit code - tells reshctl() wrapper to enable control R binding
	EnableControlRBinding = 32

	// DisableArrowKeyBindings exit code - tells reshctl() wrapper to disable arrow key bindings
	DisableArrowKeyBindings = 41
	// DisableControlRBinding exit code - tells reshctl() wrapper to disable control R binding
	DisableControlRBinding = 42
	// ReloadRcFiles exit code - tells reshctl() wrapper to reload shellrc resh file
	ReloadRcFiles = 50
	// InspectSessionHistory exit code - tells reshctl() wrapper to take current sessionID and send /inspect request to daemon
	InspectSessionHistory = 51
	// ReshStatus exit code - tells reshctl() wrapper to show RESH status (aka systemctl status)
	ReshStatus = 52
)
