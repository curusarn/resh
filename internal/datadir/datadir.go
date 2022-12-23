package datadir

import (
	"fmt"
	"os"
	"path"
)

// You should not need this caching
// It messes with proper dependency injection
// Find another way

// type dirCache struct {
// 	dir string
// 	err error
//
// 	cached bool
// }
//
// var cache dirCache
//
// func getPathNoCache() (string, error) {
// 	reshDir := "resh"
// 	xdgDir, found := os.LookupEnv("XDG_DATA_HOME")
// 	if found {
// 		return path.Join(xdgDir, reshDir), nil
// 	}
// 	homeDir, err := os.UserHomeDir()
// 	if err != nil {
// 		return "", fmt.Errorf("error while getting home dir: %w", err)
// 	}
// 	return path.Join(homeDir, ".local/share/", reshDir), nil
// }
//
// func GetPath() (string, error) {
// 	if !cache.cached {
// 		dir, err := getPathNoCache()
// 		cache = dirCache{
// 			dir:    dir,
// 			err:    err,
// 			cached: true,
// 		}
// 	}
// 	return cache.dir, cache.err
// }

func GetPath() (string, error) {
	reshDir := "resh"
	xdgDir, found := os.LookupEnv("XDG_DATA_HOME")
	if found {
		return path.Join(xdgDir, reshDir), nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error while getting home dir: %w", err)
	}
	return path.Join(homeDir, ".local/share/", reshDir), nil
}

func MakePath() (string, error) {
	path, err := GetPath()
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(path, 0755)
	// skip "exists" error
	if err != nil && !os.IsExist(err) {
		return "", fmt.Errorf("error while creating directories: %w", err)
	}
	return path, nil
}
