package recordint

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/curusarn/resh/record"
	giturls "github.com/whilp/git-urls"
	"go.uber.org/zap"
)

// SearchApp record used for sending records to RESH-CLI
type SearchApp struct {
	IsRaw     bool
	SessionID string
	DeviceID  string

	CmdLine         string
	Host            string
	Pwd             string
	Home            string // helps us to collapse /home/user to tilde
	GitOriginRemote string
	ExitCode        int

	Time float64

	// file index
	Idx int
}

func NewSearchAppFromCmdLine(cmdLine string) SearchApp {
	return SearchApp{
		IsRaw:   true,
		CmdLine: cmdLine,
	}
}

// The error handling here could be better
func NewSearchApp(sugar *zap.SugaredLogger, r *record.V1) SearchApp {
	time, err := strconv.ParseFloat(r.Time, 64)
	if err != nil {
		sugar.Errorw("Error while parsing time as float", zap.Error(err),
			"time", time)
	}
	return SearchApp{
		IsRaw:     false,
		SessionID: r.SessionID,
		CmdLine:   r.CmdLine,
		Host:      r.Device,
		Pwd:       r.Pwd,
		Home:      r.Home,
		// TODO: is this the right place to normalize the git remote
		GitOriginRemote: normalizeGitRemote(sugar, r.GitOriginRemote),
		ExitCode:        r.ExitCode,
		Time:            time,
	}
}

// TODO: maybe move this to a more appropriate place
// normalizeGitRemote helper
func normalizeGitRemote(sugar *zap.SugaredLogger, gitRemote string) string {
	gitRemote = strings.TrimSuffix(gitRemote, ".git")
	parsedURL, err := giturls.Parse(gitRemote)
	if err != nil {
		sugar.Errorw("Failed to parse git remote", zap.Error(err),
			"gitRemote", gitRemote,
		)
		return gitRemote
	}
	if parsedURL.User == nil || parsedURL.User.Username() == "" {
		parsedURL.User = url.User("git")
	}
	// TODO: figure out what scheme we want
	parsedURL.Scheme = "git+ssh"
	return parsedURL.String()
}
