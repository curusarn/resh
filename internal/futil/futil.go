package futil

import (
	"fmt"
	"io"
	"os"
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
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
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

func CreateFile(fpath string) error {
	ff, err := os.Create(fpath)
	if err != nil {
		return err
	}
	err = ff.Close()
	if err != nil {
		return err
	}
	return nil
}
