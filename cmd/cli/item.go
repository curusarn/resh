package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/curusarn/resh/pkg/records"
)

const itemLocationLenght = 30

type item struct {
	realtimeBefore float64

	// [host:]pwd
	differentHost bool
	host          string
	home          string
	samePwd       bool
	pwd           string

	// [G] [E#]
	sameGitRepo bool
	exitCode    int

	cmdLineWithColor string
	cmdLine          string

	score float64

	key string
	// cmdLineRaw string
}

type itemColumns struct {
	dateWithColor string
	date          string

	// [host:]pwd
	hostWithColor string
	host          string
	pwdTilde      string
	samePwd       bool
	//locationWithColor string
	//location          string

	// [G] [E#]
	flagsWithColor string
	flags          string

	cmdLineWithColor string
	cmdLine          string

	// score float64

	key string
	// cmdLineRaw string
}

func (i item) less(i2 item) bool {
	// reversed order
	return i.score > i2.score
}
func (i item) drawItemColumns(compactRendering bool) itemColumns {

	// DISPLAY
	// DISPLAY > date
	secs := int64(i.realtimeBefore)
	nsecs := int64((i.realtimeBefore - float64(secs)) * 1e9)
	tm := time.Unix(secs, nsecs)

	var date string
	if compactRendering {
		date = formatTimeRelativeShort(tm) + " "
	} else {
		date = formatTimeRelativeLong(tm) + " "
	}
	dateWithColor := highlightDate(date)

	// DISPLAY > location
	// DISPLAY > location > host
	host := ""
	hostWithColor := ""
	if i.differentHost {
		host += i.host + ":"
		hostWithColor += highlightHost(i.host) + ":"
	}
	// DISPLAY > location > directory
	pwdTilde := strings.Replace(i.pwd, i.home, "~", 1)

	// DISPLAY > flags
	flags := ""
	flagsWithColor := ""
	if debug {
		hitsStr := fmt.Sprintf("%.1f", i.score)
		flags += " S" + hitsStr
	}
	if i.sameGitRepo {
		flags += " G"
		flagsWithColor += " " + highlightGit("G")
	}
	if i.exitCode != 0 {
		flags += " E" + strconv.Itoa(i.exitCode)
		flagsWithColor += " " + highlightWarn("E"+strconv.Itoa(i.exitCode))
	}
	// NOTE: you can debug arbitrary metadata like this
	// flags += " <" + record.GitOriginRemote + ">"
	// flagsWithColor += " <" + record.GitOriginRemote + ">"
	return itemColumns{
		date:             date,
		dateWithColor:    dateWithColor,
		host:             host,
		hostWithColor:    hostWithColor,
		pwdTilde:         pwdTilde,
		samePwd:          i.samePwd,
		flags:            flags,
		flagsWithColor:   flagsWithColor,
		cmdLine:          i.cmdLine,
		cmdLineWithColor: i.cmdLineWithColor,
		// score:             i.score,
		key: i.key,
	}
}

func (ic itemColumns) produceLine(dateLength int, locationLength int, flagLength int, header bool, showDate bool) (string, int) {
	line := ""
	if showDate {
		date := ic.date
		for len(date) < dateLength {
			line += " "
			date += " "
		}
		// TODO: use strings.Repeat
		line += ic.dateWithColor
	}
	// LOCATION
	locationWithColor := ic.hostWithColor
	pwdLength := locationLength - len(ic.host)
	if ic.samePwd {
		locationWithColor += highlightPwd(leftCutPadString(ic.pwdTilde, pwdLength))
	} else {
		locationWithColor += leftCutPadString(ic.pwdTilde, pwdLength)
	}
	line += locationWithColor
	line += ic.flagsWithColor
	flags := ic.flags
	if flagLength < len(ic.flags) {
		log.Printf("produceLine can't specify line w/ flags shorter than the actual size. - len(flags) %v, requested %v\n", len(ic.flags), flagLength)
	}
	for len(flags) < flagLength {
		line += " "
		flags += " "
	}
	spacer := "  "
	if flagLength > 5 || header {
		// use shorter spacer
		// 		because there is likely a long flag like E130 in the view
		spacer = " "
	}
	line += spacer + ic.cmdLineWithColor

	length := dateLength + locationLength + flagLength + len(spacer) + len(ic.cmdLine)
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
	// errorExitStatus := false
	if record.ExitCode != 0 {
		// errorExitStatus = true
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

	// DISPLAY > cmdline

	// cmd := "<" + strings.ReplaceAll(record.CmdLine, "\n", ";") + ">"
	cmdLine := strings.ReplaceAll(record.CmdLine, "\n", ";")
	cmdLineWithColor := strings.ReplaceAll(cmd, "\n", ";")

	it := item{
		realtimeBefore: record.RealtimeBefore,

		differentHost: differentHost,
		host:          record.Host,
		home:          record.Home,
		samePwd:       samePwd,
		pwd:           record.Pwd,

		sameGitRepo:      sameGitRepo,
		exitCode:         record.ExitCode,
		cmdLine:          cmdLine,
		cmdLineWithColor: cmdLineWithColor,
		score:            score,
		key:              key,
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
