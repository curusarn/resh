package recordint

import (
	"net/url"
	"strconv"
	"strings"

	giturls "github.com/whilp/git-urls"
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

// NewCliRecordFromCmdLine
func NewSearchAppFromCmdLine(cmdLine string) SearchApp {
	return SearchApp{
		IsRaw:   true,
		CmdLine: cmdLine,
	}
}

// NewCliRecord from EnrichedRecord
func NewSearchApp(r *Indexed) SearchApp {
	// TODO: we used to validate records with recutil.Validate()
	// TODO: handle this error
	time, _ := strconv.ParseFloat(r.Rec.Time, 64)
	return SearchApp{
		IsRaw:     false,
		SessionID: r.Rec.SessionID,
		CmdLine:   r.Rec.CmdLine,
		Host:      r.Rec.Device,
		Pwd:       r.Rec.Pwd,
		Home:      r.Rec.Home,
		// TODO: is this the right place to normalize the git remote
		GitOriginRemote: normalizeGitRemote(r.Rec.GitOriginRemote),
		ExitCode:        r.Rec.ExitCode,
		Time:            time,

		Idx: r.Idx,
	}
}

// TODO: maybe move this to a more appropriate place
// normalizeGitRemote helper
func normalizeGitRemote(gitRemote string) string {
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
