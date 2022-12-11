package main

import (
	"github.com/curusarn/resh/cmd/control/cmd"
)

var version string
var commit string
var development string

func main() {
	cmd.Execute(version, commit, development)
}
