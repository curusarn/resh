package device

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/curusarn/resh/internal/futil"
	"github.com/google/uuid"
	isatty "github.com/mattn/go-isatty"
)

const fnameID = "device-id"
const fnameName = "device-name"

const filePerm = 0644

// Getters

func GetID(dataDir string) (string, error) {
	return readValue(dataDir, fnameID)
}

func GetName(dataDir string) (string, error) {
	return readValue(dataDir, fnameName)
}

// Install helpers

func SetupID(dataDir string) error {
	return generateIDIfUnset(dataDir)
}

func SetupName(dataDir string) error {
	return promptAndWriteNameIfUnset(dataDir)
}

func readValue(dataDir, fname string) (string, error) {
	fpath := path.Join(dataDir, fname)
	dat, err := os.ReadFile(fpath)
	if err != nil {
		return "", fmt.Errorf("could not read file with %s: %w", fname, err)
	}
	val := strings.TrimRight(string(dat), "\n")
	return val, nil
}

func generateIDIfUnset(dataDir string) error {
	fpath := path.Join(dataDir, fnameID)
	exists, err := futil.FileExists(fpath)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	rnd, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("could not get new random source: %w", err)
	}
	id := rnd.String()
	if id == "" {
		return fmt.Errorf("got invalid UUID from package")
	}
	err = os.WriteFile(fpath, []byte(id), filePerm)
	if err != nil {
		return fmt.Errorf("could not write generated ID to file: %w", err)
	}
	return nil
}

func promptAndWriteNameIfUnset(dataDir string) error {
	fpath := path.Join(dataDir, fnameName)
	exists, err := futil.FileExists(fpath)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	name, err := promptForName(fpath)
	if err != nil {
		return fmt.Errorf("error while prompting for input: %w", err)
	}
	err = os.WriteFile(fpath, []byte(name), filePerm)
	if err != nil {
		return fmt.Errorf("could not write name to file: %w", err)
	}
	return nil
}

func promptForName(fpath string) (string, error) {
	// This function should be only ran from install-utils with attached terminal
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return "", fmt.Errorf("output is not a terminal - write name of this device to '%s' to bypass this error", fpath)
	}
	host, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("could not get hostname (prompt default): %w", err)
	}
	hostStub := strings.Split(host, ".")[0]
	fmt.Printf("\nPlease choose a short name for this device (default: '%s'): ", hostStub)
	var input string
	n, err := fmt.Scanln(&input)
	if n != 1 {
		return "", fmt.Errorf("expected 1 value from prompt got %d", n)
	}
	if err != nil {
		return "", fmt.Errorf("scanln error: %w", err)
	}
	fmt.Printf("Input was: %s\n", input)
	fmt.Printf("You can change the device name at any time by editing '%s' file\n", fpath)
	return input, nil
}
