package main

import (
	"fmt"
	"os"

	"github.com/curusarn/resh/internal/datadir"
	"github.com/curusarn/resh/internal/device"
)

func setupDevice() {
	dataDir, err := datadir.MakePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to get/setup data directory: %v\n", err)
		os.Exit(1)
	}
	err = device.SetupName(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to check/setup device name: %v\n", err)
		os.Exit(1)
	}
	err = device.SetupID(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to check/setup device ID: %v\n", err)
		os.Exit(1)
	}
}
