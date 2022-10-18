package output

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

// Output wrapper for writting to logger and stdout/stderr at the same time
// useful for errors that should be presented to the user
type Output struct {
	Logger    *zap.Logger
	ErrPrefix string
}

func New(logger *zap.Logger, prefix string) *Output {
	return &Output{
		Logger:    logger,
		ErrPrefix: prefix,
	}
}

func (f *Output) Info(msg string) {
	fmt.Fprintf(os.Stdout, msg)
	f.Logger.Info(msg)
}

func (f *Output) Error(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s: %v", f.ErrPrefix, msg, err)
	f.Logger.Error(msg, zap.Error(err))
}

func (f *Output) Fatal(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s: %v", f.ErrPrefix, msg, err)
	f.Logger.Fatal(msg, zap.Error(err))
}

var msgDeamonNotRunning = `Resh-daemon didn't respond - it's probably not running.

 -> Try restarting this terminal window to bring resh-daemon back up
 -> If the problem persists you can check resh-daemon logs: ~/.local/share/resh/log.json (or ~/$XDG_DATA_HOME/resh/log.json)
 -> You can create an issue at: https://github.com/curusarn/resh/issues
`
var msgVersionMismatch = `This terminal session was started with different resh version than is installed now.
It looks like you updated resh and didn't restart this terminal.

 -> Restart this terminal window to fix that
`

func (f *Output) ErrorDaemonNotRunning(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s", f.ErrPrefix, msgDeamonNotRunning)
	f.Logger.Error("Daemon is not running", zap.Error(err))
}

func (f *Output) FatalDaemonNotRunning(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s", f.ErrPrefix, msgDeamonNotRunning)
	f.Logger.Fatal("Daemon is not running", zap.Error(err))
}

func (f *Output) ErrorVersionMismatch(installedVer, terminalVer string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n\n(installed version: %s, this terminal version: %s)",
		f.ErrPrefix, msgVersionMismatch, installedVer, terminalVer)
	f.Logger.Fatal("Version mismatch",
		zap.String("installed", installedVer),
		zap.String("terminal", terminalVer))
}

func (f *Output) FatalVersionMismatch(installedVer, terminalVer string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n\n(installed version: %s, this terminal version: %s)",
		f.ErrPrefix, msgVersionMismatch, installedVer, terminalVer)
	f.Logger.Fatal("Version mismatch",
		zap.String("installed", installedVer),
		zap.String("terminal", terminalVer))
}
