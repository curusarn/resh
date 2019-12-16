package main

import (
	"os"

	"github.com/curusarn/resh/cmd/control/cmd"
)

// version from git set during build
var version string

// commit from git set during build
var commit string

func main() {
	os.Exit(int(cmd.Execute(version, commit)))
}
