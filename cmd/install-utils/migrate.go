package main

import (
	"fmt"
	"os"

	"github.com/curusarn/resh/internal/cfg"
)

func migrateConfig() {
	err := cfg.Touch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to touch config file: %v\n", err)
		os.Exit(1)
	}
	changes, err := cfg.Migrate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to update config file: %v\n", err)
		os.Exit(1)
	}
	if changes {
		fmt.Printf("RESH config file format has changed since last update - your config was updated to reflect the changes.\n")
	}
}

func migrateHistory() {
	// homeDir, err := os.UserHomeDir()
	// if err != nil {

	// }

	// TODO: Find history in:
	//  - .resh/history.json (copy) - message user to delete the file once they confirm the new setup works
	//  - .resh_history.json (copy) - message user to delete the file once they confirm the new setup works
	//  - xdg_data/resh/history.reshjson

	// Read xdg_data/resh/history.reshjson
	// Write xdg_data/resh/history.reshjson
	// the old format can be found in the backup dir
}
