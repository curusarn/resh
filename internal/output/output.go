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

// Info outputs string to stdout and to log (as info)
// This is how we write output to users from interactive commands
// This way we have full record in logs
func (f *Output) Info(msg string) {
	fmt.Printf("%s\n", msg)
	f.Logger.Info(msg)
}

// InfoE outputs string to stdout and to log (as error)
// Passed error is only written to log
// This is how we output errors to users from interactive commands
// This way we have errors in logs
func (f *Output) InfoE(msg string, err error) {
	fmt.Printf("%s\n", msg)
	f.Logger.Error(msg, zap.Error(err))
}

// Error outputs string to stderr and to log (as error)
// This is how we output errors from non-interactive commands
func (f *Output) Error(msg string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", f.ErrPrefix, msg)
	f.Logger.Error(msg)
}

// ErrorE outputs string and error to stderr and to log (as error)
// This is how we output errors from non-interactive commands
func (f *Output) ErrorE(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s: %v\n", f.ErrPrefix, msg, err)
	f.Logger.Error(msg, zap.Error(err))
}

// FatalE outputs string and error to stderr and to log (as fatal)
// This is how we raise fatal errors from non-interactive commands
func (f *Output) FatalE(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s: %v\n", f.ErrPrefix, msg, err)
	f.Logger.Fatal(msg, zap.Error(err))
}

var msgDaemonNotRunning = `RESH daemon didn't respond - it's probably not running.

 -> Start RESH daemon manually - run: resh-daemon-start
 -> Or restart this terminal window to bring RESH daemon back up
 -> You can check logs: ~/.local/share/resh/log.json (or ~/$XDG_DATA_HOME/resh/log.json)
 -> You can create an issue at: https://github.com/curusarn/resh/issues

`
var msgTerminalVersionMismatch = `This terminal session was started with different RESH version than is installed now.
It looks like you updated RESH and didn't restart this terminal.

 -> Restart this terminal window to fix that

`

var msgDaemonVersionMismatch = `RESH daemon is running in different version than is installed now.
It looks like something went wrong during RESH update.

 -> Kill resh-daemon and then launch a new terminal window to fix that: killall resh-daemon
 -> You can create an issue at: https://github.com/curusarn/resh/issues

`

func (f *Output) InfoDaemonNotRunning(err error) {
	fmt.Printf("%s", msgDaemonNotRunning)
	f.Logger.Error("Daemon is not running", zap.Error(err))
}

func (f *Output) ErrorDaemonNotRunning(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s", f.ErrPrefix, msgDaemonNotRunning)
	f.Logger.Error("Daemon is not running", zap.Error(err))
}

func (f *Output) FatalDaemonNotRunning(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s", f.ErrPrefix, msgDaemonNotRunning)
	f.Logger.Fatal("Daemon is not running", zap.Error(err))
}

func (f *Output) InfoTerminalVersionMismatch(installedVer, terminalVer string) {
	fmt.Printf("%s(installed version: %s, this terminal version: %s)\n\n",
		msgTerminalVersionMismatch, installedVer, terminalVer)
	f.Logger.Fatal("Version mismatch",
		zap.String("installed", installedVer),
		zap.String("terminal", terminalVer))
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

func (f *Output) InfoDaemonVersionMismatch(installedVer, daemonVer string) {
	fmt.Printf("%s(installed version: %s, running daemon version: %s)\n\n",
		msgDaemonVersionMismatch, installedVer, daemonVer)
	f.Logger.Error("Version mismatch",
		zap.String("installed", installedVer),
		zap.String("daemon", daemonVer))
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
