package recutil

import (
	"errors"
	"net/url"
	"strings"

	"github.com/curusarn/resh/internal/record"
	"github.com/mattn/go-shellwords"
	giturls "github.com/whilp/git-urls"
)

// NormalizeGitRemote helper
func NormalizeGitRemote(gitRemote string) string {
	if strings.HasSuffix(gitRemote, ".git") {
		gitRemote = gitRemote[:len(gitRemote)-4]
	}
	parsedURL, err := giturls.Parse(gitRemote)
	if err != nil {
		// TODO: log this error
		return gitRemote
	}
	if parsedURL.User == nil || parsedURL.User.Username() == "" {
		parsedURL.User = url.User("git")
	}
	// TODO: figure out what scheme we want
	parsedURL.Scheme = "git+ssh"
	return parsedURL.String()
}

// Validate returns error if the record is invalid
func Validate(r *record.V1) error {
	if r.CmdLine == "" {
		return errors.New("There is no CmdLine")
	}
	if r.RealtimeBefore == 0 || r.RealtimeAfter == 0 {
		return errors.New("There is no Time")
	}
	if r.RealtimeBeforeLocal == 0 || r.RealtimeAfterLocal == 0 {
		return errors.New("There is no Local Time")
	}
	if r.RealPwd == "" || r.RealPwdAfter == "" {
		return errors.New("There is no Real Pwd")
	}
	if r.Pwd == "" || r.PwdAfter == "" {
		return errors.New("There is no Pwd")
	}
	return nil
}

// Merge two records (part1 - collect + part2 - postcollect)
func Merge(r1 *record.V1, r2 *record.V1) error {
	if r1.PartOne == false || r2.PartOne {
		return errors.New("Expected part1 and part2 of the same record - usage: Merge(part1, part2)")
	}
	if r1.SessionID != r2.SessionID {
		return errors.New("Records to merge are not from the same sesion - r1:" + r1.SessionID + " r2:" + r2.SessionID)
	}
	if r1.CmdLine != r2.CmdLine {
		return errors.New("Records to merge are not parts of the same records - r1:" + r1.CmdLine + " r2:" + r2.CmdLine)
	}
	if r1.RecordID != r2.RecordID {
		return errors.New("Records to merge do not have the same ID - r1:" + r1.RecordID + " r2:" + r2.RecordID)
	}
	r1.ExitCode = r2.ExitCode
	r1.Duration = r2.Duration

	r1.PartsMerged = true
	r1.PartOne = false
	return nil
}

// GetCommandAndFirstWord func
func GetCommandAndFirstWord(cmdLine string) (string, string, error) {
	args, err := shellwords.Parse(cmdLine)
	if err != nil {
		// Println("shellwords Error:", err, " (cmdLine: <", cmdLine, "> )")
		return "", "", err
	}
	if len(args) == 0 {
		return "", "", nil
	}
	i := 0
	for true {
		// commands in shell sometimes look like this `variable=something command argument otherArgument --option`
		// to get the command we skip over tokens that contain '='
		if strings.ContainsRune(args[i], '=') && len(args) > i+1 {
			i++
			continue
		}
		return args[i], args[0], nil
	}
	return "ERROR", "ERROR", errors.New("failed to retrieve first word of command")
}
