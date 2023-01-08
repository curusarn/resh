package output

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

// Output wrapper for writing to logger and stdout/stderr at the same time
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
	fmt.Printf("%s\n", msg)
	f.Logger.Info(msg)
}

func (f *Output) Error(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s: %v\n", f.ErrPrefix, msg, err)
	f.Logger.Error(msg, zap.Error(err))
}

func (f *Output) ErrorWOErr(msg string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", f.ErrPrefix, msg)
	f.Logger.Error(msg)
}

func (f *Output) Fatal(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s: %v\n", f.ErrPrefix, msg, err)
	f.Logger.Fatal(msg, zap.Error(err))
}

var msgDaemonNotRunning = `Resh-daemon didn't respond - it's probably not running.

 -> Try restarting this terminal window to bring resh-daemon back up
 -> If the problem persists you can check resh-daemon logs: ~/.local/share/resh/log.json (or ~/$XDG_DATA_HOME/resh/log.json)
 -> You can create an issue at: https://github.com/curusarn/resh/issues

`
var msgTerminalVersionMismatch = `This terminal session was started with different resh version than is installed now.
It looks like you updated resh and didn't restart this terminal.

 -> Restart this terminal window to fix that

`

var msgDaemonVersionMismatch = `Resh-daemon is running in different version than is installed now.
It looks like something went wrong during resh update.

 -> Kill resh-daemon and then launch a new terminal window to fix that: pkill resh-daemon
 -> You can create an issue at: https://github.com/curusarn/resh/issues

`

func (f *Output) ErrorDaemonNotRunning(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s", f.ErrPrefix, msgDaemonNotRunning)
	f.Logger.Error("Daemon is not running", zap.Error(err))
}

func (f *Output) FatalDaemonNotRunning(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s", f.ErrPrefix, msgDaemonNotRunning)
	f.Logger.Fatal("Daemon is not running", zap.Error(err))
}

func (f *Output) ErrorTerminalVersionMismatch(installedVer, terminalVer string) {
	fmt.Fprintf(os.Stderr, "%s: %s(installed version: %s, this terminal version: %s)\n\n",
		f.ErrPrefix, msgTerminalVersionMismatch, installedVer, terminalVer)
	f.Logger.Fatal("Version mismatch",
		zap.String("installed", installedVer),
		zap.String("terminal", terminalVer))
}

func (f *Output) FatalTerminalVersionMismatch(installedVer, terminalVer string) {
	fmt.Fprintf(os.Stderr, "%s: %s(installed version: %s, this terminal version: %s)\n\n",
		f.ErrPrefix, msgTerminalVersionMismatch, installedVer, terminalVer)
	f.Logger.Fatal("Version mismatch",
		zap.String("installed", installedVer),
		zap.String("terminal", terminalVer))
}

func (f *Output) ErrorDaemonVersionMismatch(installedVer, daemonVer string) {
	fmt.Fprintf(os.Stderr, "%s: %s(installed version: %s, running daemon version: %s)\n\n",
		f.ErrPrefix, msgDaemonVersionMismatch, installedVer, daemonVer)
	f.Logger.Error("Version mismatch",
		zap.String("installed", installedVer),
		zap.String("daemon", daemonVer))
}

func (f *Output) FatalDaemonVersionMismatch(installedVer, daemonVer string) {
	fmt.Fprintf(os.Stderr, "%s: %s(installed version: %s, running daemon version: %s)\n\n",
		f.ErrPrefix, msgDaemonVersionMismatch, installedVer, daemonVer)
	f.Logger.Fatal("Version mismatch",
		zap.String("installed", installedVer),
		zap.String("daemon", daemonVer))
}
