package main

import (
	"fmt"
	"os"
	"path"

	"github.com/curusarn/resh/internal/cfg"
	"github.com/curusarn/resh/internal/datadir"
	"github.com/curusarn/resh/internal/futil"
	"github.com/curusarn/resh/internal/output"
	"github.com/curusarn/resh/internal/recio"
)

func migrateConfig(out *output.Output) {
	err := cfg.Touch()
	if err != nil {
		out.Fatal("ERROR: Failed to touch config file", err)
	}
	changes, err := cfg.Migrate()
	if err != nil {
		out.Fatal("ERROR: Failed to update config file", err)
	}
	if changes {
		out.Info("RESH config file format has changed since last update - your config was updated to reflect the changes.")
	}
}

func migrateHistory(out *output.Output) {
	migrateHistoryLocation(out)
	migrateHistoryFormat(out)
}

// find first existing history and use it
// don't bother with merging of history in multiple locations - it could get messy and it shouldn't be necessary
func migrateHistoryLocation(out *output.Output) {
	dataDir, err := datadir.MakePath()
	if err != nil {
		out.Fatal("ERROR: Failed to get data directory", err)
	}
	// TODO: de-hardcode this
	historyPath := path.Join(dataDir, "resh/history.reshjson")

	exists, err := futil.FileExists(historyPath)
	if err != nil {
		out.Fatal("ERROR: Failed to check history file", err)
	}
	if exists {
		out.Info(fmt.Sprintf("Found history file in '%s'", historyPath))
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		out.Fatal("ERROR: Failed to get user home directory", err)
	}

	legacyHistoryPaths := []string{
		path.Join(homeDir, ".resh_history.json"),
		path.Join(homeDir, ".resh/history.json"),
	}
	for _, path := range legacyHistoryPaths {
		exists, err = futil.FileExists(path)
		if err != nil {
			out.Fatal("ERROR: Failed to check legacy history file", err)
		}
		if exists {
			out.Info(fmt.Sprintf("Copying history file to new location: '%s' -> '%s' ...", path, historyPath))
			err = futil.CopyFile(path, historyPath)
			if err != nil {
				out.Fatal("ERROR: Failed to copy history file", err)
			}
			out.Info("History file copied successfully")
			return
		}
	}
}

func migrateHistoryFormat(out *output.Output) {
	dataDir, err := datadir.MakePath()
	if err != nil {
		out.Fatal("ERROR: Could not get user data directory", err)
	}
	// TODO: de-hardcode this
	historyPath := path.Join(dataDir, "history.reshjson")
	historyPathBak := historyPath + ".bak"

	exists, err := futil.FileExists(historyPath)
	if err != nil {
		out.Fatal("ERROR: Failed to check existence of history file", err)
	}
	if !exists {
		out.ErrorWOErr("There is no history file - this is normal if you are installing RESH for the first time on this device")
		err = futil.TouchFile(historyPath)
		if err != nil {
			out.Fatal("ERROR: Failed to touch history file", err)
		}
		os.Exit(0)
	}

	err = futil.CopyFile(historyPath, historyPathBak)
	if err != nil {
		out.Fatal("ERROR: Could not back up history file", err)
	}

	rio := recio.New(out.Logger.Sugar())

	recs, err := rio.ReadAndFixFile(historyPath, 3)
	if err != nil {
		out.Fatal("ERROR: Could not load history file", err)
	}
	err = rio.OverwriteFile(historyPath, recs)
	if err != nil {
		out.Error("ERROR: Could not update format of history file", err)

		err = futil.CopyFile(historyPathBak, historyPath)
		if err != nil {
			out.Fatal("ERROR: Could not restore history file from backup!", err)
			// TODO: history restoration tutorial
		}
		out.Info("History file was restored to the original form")
	}
}
