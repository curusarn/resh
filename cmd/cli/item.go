package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/curusarn/resh/pkg/records"
)

type item struct {
	// dateWithColor string
	// date          string

	// [host:]pwd
	locationWithColor string
	location          string

	// [G] [E#]
	flagsWithColor string
	flags          string

	cmdLineWithColor string
	cmdLine          string

	score float64

	key string
	// cmdLineRaw string
}

func (i item) less(i2 item) bool {
	// reversed order
	return i.score > i2.score
}

func (i item) produceLine(flagLength int) (string, int) {
	line := ""
	line += i.locationWithColor
	line += i.flagsWithColor
	flags := i.flags
	if flagLength < len(i.flags) {
		log.Printf("produceLine can't specify line w/ flags shorter than the actual size. - len(flags) %v, requested %v\n", len(i.flags), flagLength)
	}
	for len(flags) < flagLength {
		line += " "
		flags += " "
	}
	spacer := "  "
	if flagLength > 5 {
		// use shorter spacer
		// 		because there is likely a long flag like E130 in the view
		spacer = " "
	}
	line += spacer + i.cmdLineWithColor

	length := len(i.location) + flagLength + len(spacer) + len(i.cmdLine)
	return line, length
}

func leftCutPadString(str string, newLen int) string {
	dots := "…"
	strLen := len(str)
	if newLen > strLen {
		return strings.Repeat(" ", newLen-strLen) + str
	} else if newLen < strLen {
		return dots + str[strLen-newLen+1:]
	}
	return str
}

func rightCutPadString(str string, newLen int) string {
	dots := "…"
	strLen := len(str)
	if newLen > strLen {
		return str + strings.Repeat(" ", newLen-strLen)
	} else if newLen < strLen {
		return str[:newLen-1] + dots
	}
	return str
}

// proper match for path is when whole directory is matched
// proper match for command is when term matches word delimeted by whitespace
func properMatch(str, term, padChar string) bool {
	if strings.Contains(padChar+str+padChar, padChar+term+padChar) {
		return true
	}
	return false
}

// newItemFromRecordForQuery creates new item from record based on given query
//		returns error if the query doesn't match the record
func newItemFromRecordForQuery(record records.CliRecord, query query, debug bool) (item, error) {
	const hitScore = 1.0
	const hitScoreConsecutive = 0.1
	const properMatchScore = 0.3
	const actualPwdScore = 0.9
	const nonZeroExitCodeScorePenalty = 0.5
	const sameGitRepoScore = 0.7
	// const sameGitRepoScoreExtra = 0.0
	const differentHostScorePenalty = 0.2

	// nonZeroExitCodeScorePenalty + differentHostScorePenalty

	score := 0.0
	anyHit := false
	cmd := record.CmdLine
	for _, term := range query.terms {
		termHit := false
		if strings.Contains(record.CmdLine, term) {
			anyHit = true
			if termHit == false {
				score += hitScore
			} else {
				score += hitScoreConsecutive
			}
			termHit = true
			if properMatch(cmd, term, " ") {
				score += properMatchScore
			}
			cmd = strings.ReplaceAll(cmd, term, highlightMatch(term))
			// NO continue
		}
	}
	// actual pwd matches
	// N terms can only produce:
	//		-> N matches against the command
	//		-> 1 extra match for the actual directory match
	sameGitRepo := false
	if query.gitOriginRemote != "" && query.gitOriginRemote == record.GitOriginRemote {
		sameGitRepo = true
	}

	samePwd := false
	if record.Pwd == query.pwd {
		anyHit = true
		samePwd = true
		score += actualPwdScore
	} else if sameGitRepo {
		anyHit = true
		score += sameGitRepoScore
	}

	differentHost := false
	if record.Host != query.host {
		differentHost = true
		score -= differentHostScorePenalty
	}
	errorExitStatus := false
	if record.ExitCode != 0 {
		errorExitStatus = true
		score -= nonZeroExitCodeScorePenalty
	}
	if score <= 0 && !anyHit {
		return item{}, errors.New("no match for given record and query")
	}

	// KEY for deduplication

	unlikelySeparator := "|||||"
	key := record.CmdLine + unlikelySeparator + record.Pwd + unlikelySeparator +
		record.GitOriginRemote + unlikelySeparator + record.Host
	// + strconv.Itoa(record.ExitCode) + unlikelySeparator

	// DISPLAY
	// DISPLAY > date
	// TODO

	// DISPLAY > location
	location := ""
	locationWithColor := ""
	if differentHost {
		location += record.Host + ":"
		locationWithColor += highlightHost(record.Host) + ":"
	}
	const locationLenght = 30
	// const locationLenght = 20 // small screenshots
	pwdLength := locationLenght - len(location)
	pwdTilde := strings.Replace(record.Pwd, record.Home, "~", 1)
	location += leftCutPadString(pwdTilde, pwdLength)
	if samePwd {
		locationWithColor += highlightPwd(leftCutPadString(pwdTilde, pwdLength))
	} else {
		locationWithColor += leftCutPadString(pwdTilde, pwdLength)
	}

	// DISPLAY > flags
	flags := ""
	flagsWithColor := ""
	if debug {
		hitsStr := fmt.Sprintf("%.1f", score)
		flags += " S" + hitsStr
	}
	if sameGitRepo {
		flags += " G"
		flagsWithColor += " " + highlightGit("G")
	}
	if errorExitStatus {
		flags += " E" + strconv.Itoa(record.ExitCode)
		flagsWithColor += " " + highlightWarn("E"+strconv.Itoa(record.ExitCode))
	}
	// NOTE: you can debug arbitrary metadata like this
	// flags += " <" + record.GitOriginRemote + ">"
	// flagsWithColor += " <" + record.GitOriginRemote + ">"

	// DISPLAY > cmdline

	// cmd := "<" + strings.ReplaceAll(record.CmdLine, "\n", ";") + ">"
	cmdLine := strings.ReplaceAll(record.CmdLine, "\n", ";")
	cmdLineWithColor := strings.ReplaceAll(cmd, "\n", ";")

	it := item{
		location:          location,
		locationWithColor: locationWithColor,
		flags:             flags,
		flagsWithColor:    flagsWithColor,
		cmdLine:           cmdLine,
		cmdLineWithColor:  cmdLineWithColor,
		score:             score,
		key:               key,
	}
	return it, nil
}

type rawItem struct {
	cmdLineWithColor string
	cmdLine          string

	hits float64

	key string
	// cmdLineRaw string
}

// newRawItemFromRecordForQuery creates new item from record based on given query
//		returns error if the query doesn't match the record
func newRawItemFromRecordForQuery(record records.CliRecord, terms []string, debug bool) (rawItem, error) {
	const hitScore = 1.0
	const hitScoreConsecutive = 0.1
	const properMatchScore = 0.3

	hits := 0.0
	anyHit := false
	cmd := record.CmdLine
	for _, term := range terms {
		termHit := false
		if strings.Contains(record.CmdLine, term) {
			anyHit = true
			if termHit == false {
				hits += hitScore
			} else {
				hits += hitScoreConsecutive
			}
			termHit = true
			if properMatch(cmd, term, " ") {
				hits += properMatchScore
			}
			cmd = strings.ReplaceAll(cmd, term, highlightMatch(term))
			// NO continue
		}
	}
	_ = anyHit
	// KEY for deduplication
	key := record.CmdLine

	// DISPLAY > cmdline

	// cmd := "<" + strings.ReplaceAll(record.CmdLine, "\n", ";") + ">"
	cmdLine := strings.ReplaceAll(record.CmdLine, "\n", ";")
	cmdLineWithColor := strings.ReplaceAll(cmd, "\n", ";")

	it := rawItem{
		cmdLine:          cmdLine,
		cmdLineWithColor: cmdLineWithColor,
		hits:             hits,
		key:              key,
	}
	return it, nil
}
