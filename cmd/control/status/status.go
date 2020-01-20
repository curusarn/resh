package status

// Code - exit code of the resh-control command
type Code int

const (
	// Success exit code
	Success Code = 0
	// Fail exit code
	Fail = 1
	// EnableResh exit code -  tells reshctl() wrapper to enable resh
	// EnableResh = 100

	// EnableArrowKeyBindings exit code - tells reshctl() wrapper to enable arrow key bindings
	EnableArrowKeyBindings = 101
	// EnableControlRBinding exit code - tells reshctl() wrapper to enable control R binding
	EnableControlRBinding = 102
	// DisableResh exit code -  tells reshctl() wrapper to enable resh
	// DisableResh = 110

	// DisableArrowKeyBindings exit code - tells reshctl() wrapper to disable arrow key bindings
	DisableArrowKeyBindings = 111
	// DisableControlRBinding exit code - tells reshctl() wrapper to disable control R binding
	DisableControlRBinding = 112
	// ReloadRcFiles exit code - tells reshctl() wrapper to reload shellrc resh file
	ReloadRcFiles = 200
	// InspectSessionHistory exit code - tells reshctl() wrapper to take current sessionID and send /inspect request to daemon
	InspectSessionHistory = 201
	// ReshStatus exit code - tells reshctl() wrapper to show RESH status (aka systemctl status)
	ReshStatus = 202
)
