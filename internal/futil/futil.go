// futil implements common file-related utilities
package futil

import (
	"fmt"
	"io"
	"os"
	"time"
)

func CopyFile(source, dest string) error {
	from, err := os.Open(source)
	if err != nil {
		return err
	}
	defer from.Close()

	// This is equivalent to: os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
	to, err := os.Create(dest)
	if err != nil {
		return err
	}

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return to.Close()
}

func FileExists(fpath string) (bool, error) {
	_, err := os.Stat(fpath)
	if err == nil {
		// File exists
		return true, nil
	}
	if os.IsNotExist(err) {
		// File doesn't exist
		return false, nil
	}
	// Any other error
	return false, fmt.Errorf("could not stat file: %w", err)
}

// TouchFile touches file
// Returns true if file was created false otherwise
func TouchFile(fpath string) (bool, error) {
	exists, err := FileExists(fpath)
	if err != nil {
		return false, err
	}

	file, err := os.OpenFile(fpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return false, fmt.Errorf("could not open/create file: %w", err)
	}
	err = file.Close()
	if err != nil {
		return false, fmt.Errorf("could not close file: %w", err)
	}
	return !exists, nil
}

func getBackupPath(fpath string) string {
	ext := fmt.Sprintf(".backup-%d", time.Now().Unix())
	return fpath + ext
}

// BackupFile backups file using unique suffix
// Returns path to backup
func BackupFile(fpath string) (*RestorableFile, error) {
	fpathBackup := getBackupPath(fpath)
	exists, err := FileExists(fpathBackup)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("backup already exists in the determined path")
	}
	err = CopyFile(fpath, fpathBackup)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file: %w ", err)
	}
	rf := RestorableFile{
		Path:       fpath,
		PathBackup: fpathBackup,
	}
	return &rf, nil
}

type RestorableFile struct {
	Path       string
	PathBackup string
}

func (r RestorableFile) Restore() error {
	return restoreFileFromBackup(r.Path, r.PathBackup)
}

func restoreFileFromBackup(fpath, fpathBak string) error {
	exists, err := FileExists(fpathBak)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("backup not found in given path: no such file or directory: %s", fpathBak)
	}
	err = CopyFile(fpathBak, fpath)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w ", err)
	}
	return nil
}
