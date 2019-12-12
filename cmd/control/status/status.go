package status

// Code - exit code of the resh-control command
type Code int

const (
	// Success exit code
	Success Code = 0
	// Fail exit code
	Fail = 1
	// EnableArrowKeyBindings exit code - tells reshctl() wrapper to enable arrow key bindings
	EnableArrowKeyBindings = 101
	// DisableArrowKeyBindings exit code - tells reshctl() wrapper to disable arrow key bindings
	DisableArrowKeyBindings = 111
	// ReloadRcFiles exit code - tells reshctl() wrapper to reload shellrc resh file
	ReloadRcFiles = 200
)
