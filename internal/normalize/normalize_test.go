package normalize_test

import (
	"testing"

	"github.com/curusarn/resh/internal/normalize"
	"go.uber.org/zap"
)

// TestLeftCutPadString
func TestGitRemote(t *testing.T) {
	sugar := zap.NewNop().Sugar()

	data := [][]string{
		{
			"git@github.com:curusarn/resh.git",           // git
			"git@github.com:curusarn/resh",               // git no ".git"
			"http://github.com/curusarn/resh.git",        // http
			"https://github.com/curusarn/resh.git",       // https
			"ssh://git@github.com/curusarn/resh.git",     // ssh
			"git+ssh://git@github.com/curusarn/resh.git", // git+ssh
		},
		{
			"git@host.example.com:org/user/repo.git",           // git
			"git@host.example.com:org/user/repo",               // git no ".git"
			"http://host.example.com/org/user/repo.git",        // http
			"https://host.example.com/org/user/repo.git",       // https
			"ssh://git@host.example.com/org/user/repo.git",     // ssh
			"git+ssh://git@host.example.com/org/user/repo.git", // git+ssh
		},
	}

	for _, arr := range data {
		n := len(arr)
		for i := 0; i < n-1; i++ {
			for j := i + 1; j < n; j++ {
				one := normalize.GitRemote(sugar, arr[i])
				two := normalize.GitRemote(sugar, arr[j])
				if one != two {
					t.Fatalf("Normalized git remotes should match for '%s' and '%s'\n -> got '%s' != '%s'",
						arr[i], arr[j], one, two)
				}
			}
		}
	}

	empty := normalize.GitRemote(sugar, "")
	if len(empty) != 0 {
		t.Fatalf("Normalized git remotes for '' should be ''\n -> got '%s'", empty)
	}
}
