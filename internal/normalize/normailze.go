package normalize

import (
	"net/url"
	"strings"

	giturls "github.com/whilp/git-urls"
	"go.uber.org/zap"
)

// GitRemote helper
// Returns normalized git remote - valid even on error
func GitRemote(sugar *zap.SugaredLogger, gitRemote string) string {
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
