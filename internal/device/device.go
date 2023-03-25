// device implements helpers that get/set device config files
package device

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/curusarn/resh/internal/futil"
	"github.com/curusarn/resh/internal/output"
	"github.com/google/uuid"
	isatty "github.com/mattn/go-isatty"
)

const fnameID = "device-id"
const fnameName = "device-name"

const fpathIDLegacy = ".resh/resh-uuid"

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
	return setIDIfUnset(dataDir)
}

func SetupName(out *output.Output, dataDir string) error {
	return promptAndWriteNameIfUnset(out, dataDir)
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

func setIDIfUnset(dataDir string) error {
	fpath := path.Join(dataDir, fnameID)
	exists, err := futil.FileExists(fpath)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	// Try copy device ID from legacy location
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user home: %w", err)
	}
	fpathLegacy := path.Join(homeDir, fpathIDLegacy)
	exists, err = futil.FileExists(fpath)
	if err != nil {
		return err
	}
	if exists {
		futil.CopyFile(fpathLegacy, fpath)
		if err != nil {
			return fmt.Errorf("could not copy device ID from legacy location: %w", err)
		}
		return nil
	}

	// Generate new device ID
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

func promptAndWriteNameIfUnset(out *output.Output, dataDir string) error {
	fpath := path.Join(dataDir, fnameName)
	exists, err := futil.FileExists(fpath)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	name, err := promptForName(out, fpath)
	if err != nil {
		return fmt.Errorf("error while prompting for input: %w", err)
	}
	err = os.WriteFile(fpath, []byte(name), filePerm)
	if err != nil {
		return fmt.Errorf("could not write name to file: %w", err)
	}
	return nil
}

func promptForName(out *output.Output, fpath string) (string, error) {
	// This function should be only ran from install-utils with attached terminal
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		return "", fmt.Errorf("output is not a terminal - write name of this device to '%s' to bypass this error", fpath)
	}
	host, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("could not get hostname (prompt default): %w", err)
	}
	hostStub := strings.Split(host, ".")[0]
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nChoose a short name for this device (default: '%s'): ", hostStub)
	input, err := reader.ReadString('\n')
	name := strings.TrimRight(input, "\n")
	if err != nil {
		return "", fmt.Errorf("reader error: %w", err)
	}
	if name == "" {
		out.Info("Got no input - using default ...")
		name = hostStub
	}
	out.Info(fmt.Sprintf("Device name set to '%s'", name))
	fmt.Printf("You can change the device name at any time by editing '%s' file\n", fpath)
	return input, nil
}
