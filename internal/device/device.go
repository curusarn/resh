package device

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func GetID(dataDir string) (string, error) {
	fname := "device-id"
	dat, err := os.ReadFile(path.Join(dataDir, fname))
	if err != nil {
		return "", fmt.Errorf("could not read file with device-id: %w", err)
	}
	id := strings.TrimRight(string(dat), "\n")
	return id, nil
}

func GetName(dataDir string) (string, error) {
	fname := "device-name"
	dat, err := os.ReadFile(path.Join(dataDir, fname))
	if err != nil {
		return "", fmt.Errorf("could not read file with device-name: %w", err)
	}
	name := strings.TrimRight(string(dat), "\n")
	return name, nil
}

// TODO: implement, possibly with a better name
// func CheckID(dataDir string) (string, error) {
// 	fname := "device-id"
// 	dat, err := os.ReadFile(path.Join(dataDir, fname))
// 	if err != nil {
// 		return "", fmt.Errorf("could not read file with device-id: %w", err)
// 	}
// 	id := strings.TrimRight(string(dat), "\n")
// 	return id, nil
// }
//
// func CheckName(dataDir string) (string, error) {
// 	fname := "device-id"
// 	dat, err := os.ReadFile(path.Join(dataDir, fname))
// 	if err != nil {
// 		return "", fmt.Errorf("could not read file with device-id: %w", err)
// 	}
// 	id := strings.TrimRight(string(dat), "\n")
// 	return id, nil
// }
