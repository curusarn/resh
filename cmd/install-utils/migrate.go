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

func printRecoveryInfo(rf *futil.RestorableFile) {
	fmt.Printf(" -> Backup is '%s'\n"+
		" -> Original file location is '%s'\n"+
		" -> Please copy the backup over the file - run: cp -f '%s' '%s'\n\n",
		rf.PathBackup, rf.Path,
		rf.PathBackup, rf.Path,
	)
}

func migrateAll(out *output.Output) {
	cfgBackup, err := migrateConfig(out)
	if err != nil {
		// out.InfoE("Failed to update config file format", err)
		out.FatalE("Failed to update config file format", err)
	}
	err = migrateHistory(out)
	if err != nil {
		errHist := err
		out.InfoE("Failed to update RESH history", errHist)
		out.Info("Restoring config from backup ...")
		err = cfgBackup.Restore()
		if err != nil {
			out.InfoE("FAILED TO RESTORE CONFIG FROM BACKUP!", err)
			printRecoveryInfo(cfgBackup)
		} else {
			out.Info("Config file was restored successfully")
		}
		out.FatalE("Failed to update history", errHist)
	}
}

func migrateConfig(out *output.Output) (*futil.RestorableFile, error) {
	cfgPath, err := cfg.GetPath()
	if err != nil {
		return nil, fmt.Errorf("could not get config file path: %w", err)
	}

	// Touch config to get rid of edge-cases
	created, err := futil.TouchFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to touch config file: %w", err)
	}

	// Backup
	backup, err := futil.BackupFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("could not backup config file: %w", err)
	}

	// Migrate
	changes, err := cfg.Migrate()
	if err != nil {
		// Restore
		errMigrate := err
		errMigrateWrap := fmt.Errorf("failed to update config file: %w", errMigrate)
		out.InfoE("Failed to update config file format", errMigrate)
		out.Info("Restoring config from backup ...")
		err = backup.Restore()
		if err != nil {
			out.InfoE("FAILED TO RESTORE CONFIG FROM BACKUP!", err)
			printRecoveryInfo(backup)
		} else {
			out.Info("Config file was restored successfully")
		}
		// We are returning the root cause - there might be a better solution how to report the errors
		return nil, errMigrateWrap
	}
	if created {
		out.Info(fmt.Sprintf("RESH config created in '%s'", cfgPath))
	} else if changes {
		out.Info("RESH config file format has changed since last update - your config was updated to reflect the changes.")
	}
	return backup, nil
}

func migrateHistory(out *output.Output) error {
	err := migrateHistoryLocation(out)
	if err != nil {
		return fmt.Errorf("failed to move history to new location %w", err)
	}
	return migrateHistoryFormat(out)
}

// Find first existing history and use it
// Don't bother with merging of history in multiple locations - it could get messy and it shouldn't be necessary
func migrateHistoryLocation(out *output.Output) error {
	dataDir, err := datadir.MakePath()
	if err != nil {
		return fmt.Errorf("failed to get data directory: %w", err)
	}
	historyPath := path.Join(dataDir, datadir.HistoryFileName)

	exists, err := futil.FileExists(historyPath)
	if err != nil {
		return fmt.Errorf("failed to check history file: %w", err)
	}
	if exists {
		// TODO: get rid of this output (later)
		out.Info(fmt.Sprintf("Found history file in '%s' - nothing to move", historyPath))
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	legacyHistoryPaths := []string{
		path.Join(homeDir, ".resh_history.json"),
		path.Join(homeDir, ".resh/history.json"),
	}
	for _, path := range legacyHistoryPaths {
		exists, err = futil.FileExists(path)
		if err != nil {
			return fmt.Errorf("failed to check existence of legacy history file: %w", err)
		}
		if exists {
			// TODO: maybe get rid of this output later
			out.Info(fmt.Sprintf("Copying history file to new location: '%s' -> '%s' ...", path, historyPath))
			err = futil.CopyFile(path, historyPath)
			if err != nil {
				return fmt.Errorf("failed to copy history file: %w", err)
			}
			out.Info("History file copied successfully")
			return nil
		}
	}
	// out.Info("WARNING: No RESH history file found (this is normal during new installation)")
	return nil
}

func migrateHistoryFormat(out *output.Output) error {
	dataDir, err := datadir.MakePath()
	if err != nil {
		return fmt.Errorf("could not get user data directory: %w", err)
	}
	historyPath := path.Join(dataDir, datadir.HistoryFileName)

	exists, err := futil.FileExists(historyPath)
	if err != nil {
		return fmt.Errorf("failed to check existence of history file: %w", err)
	}
	if !exists {
		out.Error("There is no RESH history file - this is normal if you are installing RESH for the first time on this device")
		_, err = futil.TouchFile(historyPath)
		if err != nil {
			return fmt.Errorf("failed to touch history file: %w", err)
		}
		return nil
	}

	backup, err := futil.BackupFile(historyPath)
	if err != nil {
		return fmt.Errorf("could not back up history file: %w", err)
	}

	rio := recio.New(out.Logger.Sugar())

	recs, err := rio.ReadAndFixFile(historyPath, 3)
	if err != nil {
		return fmt.Errorf("could not load history file: %w", err)
	}
	err = rio.OverwriteFile(historyPath, recs)
	if err != nil {
		// Restore
		errMigrate := err
		errMigrateWrap := fmt.Errorf("failed to update format of history file: %w", errMigrate)
		out.InfoE("Failed to update RESH history file format", errMigrate)
		out.Info("Restoring RESH history from backup ...")
		err = backup.Restore()
		if err != nil {
			out.InfoE("FAILED TO RESTORE RESH HISTORY FROM BACKUP!", err)
			printRecoveryInfo(backup)
		} else {
			out.Info("RESH history file was restored successfully")
		}
		// We are returning the root cause - there might be a better solution how to report the errors
		return errMigrateWrap
	}
	return nil
}
