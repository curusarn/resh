package main

import (
	"fmt"
	"os"
)

// info passed during build
var version string
var commit string
var developement bool

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "ERROR: Not eonugh arguments\n")
		printUsage(os.Stderr)
	}
	command := os.Args[1]
	switch command {
	case "backup":
		backup()
	case "rollback":
		rollback()
	case "migrate-config":
		migrateConfig()
	case "migrate-history":
		migrateHistory()
	case "help":
		printUsage(os.Stdout)
	default:
		fmt.Fprintf(os.Stderr, "ERROR: Unknown command: %s\n", command)
		printUsage(os.Stderr)
	}
}

func printUsage(f *os.File) {
	usage := `
USAGE: ./install-utils COMMAND
Utils used during RESH instalation.

COMMANDS:
  backup		backup resh installation and data
  rollback		restore resh installation and data from backup
  migrate-config	update config to reflect updates
  migrate-history	update history to reflect updates
  help			show this help

`
	fmt.Fprintf(f, usage)
}
