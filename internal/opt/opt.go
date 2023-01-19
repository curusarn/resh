package opt

import (
	"fmt"
	"os"

	"github.com/curusarn/resh/internal/output"
)

// HandleVersionOpts reads the first option and handles it
// This is a helper for resh-{collect,postcollect,session-init} commands
func HandleVersionOpts(out *output.Output, args []string, version, commit string) []string {
	if len(os.Args) == 0 {
		return os.Args[1:]
	}
	// We use go-like options because of backwards compatibility.
	// Not ideal but we should support them because they have worked once
	// and adding "more correct" variants would mean supporting more variants.
	switch os.Args[1] {
	case "-version":
		fmt.Print(version)
		os.Exit(0)
	case "-revision":
		fmt.Print(commit)
		os.Exit(0)
	case "-requireVersion":
		if len(os.Args) < 3 {
			out.FatalTerminalVersionMismatch(version, "")
		}
		if os.Args[2] != version {
			out.FatalTerminalVersionMismatch(version, os.Args[2])
		}
		return os.Args[3:]
	}
	return os.Args[1:]
}
