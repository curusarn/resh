package main

import (
	"flag"
	"fmt"
	"os"
)

// info passed during build
var version string
var commit string
var developement bool

func main() {
	var command string
	flag.StringVar(&command, "command", "", "Utility to run")
	flag.Parse()

	switch command {
	case "backup":
		backup()
	case "rollback":
		// FIXME
		panic("Rollback not implemented yet!")
		// rollback()
	case "migrate-config":
		migrateConfig()
	case "migrate-history":
		migrateHistory()
	case "help":
		printUsage(os.Stdout)
	default:
		fmt.Fprintf(os.Stderr, "ERROR: Unknown command")
		printUsage(os.Stderr)
	}
}

func printUsage(f *os.File) {
	usage := `
	Utils used during resh instalation	

	USAGE: ./install-utils COMMAND
	COMMANDS:
	  backup		backup resh installation and data
	  rollback		restore resh installation and data from backup
	  migrate-config	update config to reflect updates
	  migrate-history	update history to reflect updates
	  help			show this help
	`
	fmt.Fprintf(f, usage)
}
