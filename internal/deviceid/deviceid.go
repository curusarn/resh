package deviceid

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func Get(dataDir string) (string, error) {
	fname := "device-id"
	dat, err := os.ReadFile(path.Join(dataDir, fname))
	if err != nil {
		return "", fmt.Errorf("could not read file with device-id: %w", err)
	}
	id := strings.TrimRight(string(dat), "\n")
	return id, nil
}
