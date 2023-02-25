package searchapp

import (
	"sort"
	"strings"

	"github.com/curusarn/resh/internal/normalize"
	"go.uber.org/zap"
)

// Query holds information that is used for result scoring
type Query struct {
	terms           []string
	host            string
	pwd             string
	gitOriginRemote string
	// pwdTilde string
}

func isValidTerm(term string) bool {
	if len(term) == 0 {
		return false
	}
	if strings.Contains(term, " ") {
		return false
	}
	return true
}

func filterTerms(terms []string) []string {
	var newTerms []string
	for _, term := range terms {
		if isValidTerm(term) {
			newTerms = append(newTerms, term)
		}
	}
	return newTerms
}

// NewQueryFromString .
func NewQueryFromString(sugar *zap.SugaredLogger, queryInput string, host string, pwd string, gitOriginRemote string, debug bool) Query {
	terms := strings.Fields(queryInput)
	var logStr string
	for _, term := range terms {
		logStr += " <" + term + ">"
	}
	terms = filterTerms(terms)
	logStr = ""
	for _, term := range terms {
		logStr += " <" + term + ">"
	}
	sort.SliceStable(terms, func(i, j int) bool { return len(terms[i]) < len(terms[j]) })
	return Query{
		terms:           terms,
		host:            host,
		pwd:             pwd,
		gitOriginRemote: normalize.GitRemote(sugar, gitOriginRemote),
	}
}

// GetRawTermsFromString .
func GetRawTermsFromString(queryInput string, debug bool) []string {
	terms := strings.Fields(queryInput)
	var logStr string
	for _, term := range terms {
		logStr += " <" + term + ">"
	}
	terms = filterTerms(terms)
	logStr = ""
	for _, term := range terms {
		logStr += " <" + term + ">"
	}
	return terms
}
