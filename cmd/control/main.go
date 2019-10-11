package main

import (
	"os"

	"github.com/curusarn/resh/cmd/control/cmd"
)

// Version from git set during build
var Version string

// Revision from git set during build
var Revision string

func main() {
	os.Exit(int(cmd.Execute()))
}
