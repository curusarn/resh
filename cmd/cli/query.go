package main

import (
	"log"
	"strings"
)

type query struct {
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

func newQueryFromString(queryInput string, host string, pwd string, gitOriginRemote string) query {
	if debug {
		log.Println("QUERY input = <" + queryInput + ">")
	}
	terms := strings.Fields(queryInput)
	var logStr string
	for _, term := range terms {
		logStr += " <" + term + ">"
	}
	if debug {
		log.Println("QUERY raw terms =" + logStr)
	}
	terms = filterTerms(terms)
	logStr = ""
	for _, term := range terms {
		logStr += " <" + term + ">"
	}
	if debug {
		log.Println("QUERY filtered terms =" + logStr)
		log.Println("QUERY pwd =" + pwd)
	}
	return query{
		terms:           terms,
		host:            host,
		pwd:             pwd,
		gitOriginRemote: gitOriginRemote,
	}
}

func getRawTermsFromString(queryInput string) []string {
	if debug {
		log.Println("QUERY input = <" + queryInput + ">")
	}
	terms := strings.Fields(queryInput)
	var logStr string
	for _, term := range terms {
		logStr += " <" + term + ">"
	}
	if debug {
		log.Println("QUERY raw terms =" + logStr)
	}
	terms = filterTerms(terms)
	logStr = ""
	for _, term := range terms {
		logStr += " <" + term + ">"
	}
	return terms
}
