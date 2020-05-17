package searchapp

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/curusarn/resh/pkg/records"
)

const itemLocationLenght = 30

// Item holds item info for normal mode
type Item struct {
	isRaw bool

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

	CmdLineWithColor string
	CmdLine          string

	Score float64

	Key string
	// cmdLineRaw string
}

// ItemColumns holds rendered columns
type ItemColumns struct {
	DateWithColor string
	Date          string

	// [host:]pwd
	HostWithColor string
	Host          string
	PwdTilde      string
	samePwd       bool
	//locationWithColor string
	//location          string

	// [G] [E#]
	FlagsWithColor string
	Flags          string

	CmdLineWithColor string
	CmdLine          string

	// score float64

	Key string
	// cmdLineRaw string
}

func (i Item) less(i2 Item) bool {
	// reversed order
	return i.Score > i2.Score
}

func splitStatusLineToLines(statusLine string, printedLineLength, realLineLength int) []string {
	var statusLineSlice []string
	// status line
	var idxSt, idxEnd int
	var nextLine bool
	tab := "    "
	tabSize := len(tab)
	for idxSt < len(statusLine) {
		idxEnd = idxSt + printedLineLength
		if nextLine {
			idxEnd -= tabSize
		}

		if idxEnd > len(statusLine) {
			idxEnd = len(statusLine)
		}
		str := statusLine[idxSt:idxEnd]

		indent := " "
		if nextLine {
			indent += tab
		}
		statusLineSlice = append(statusLineSlice, highlightStatus(rightCutPadString(indent+str, realLineLength))+"\n")
		idxSt += printedLineLength
		nextLine = true
	}
	return statusLineSlice
}

// DrawStatusLine ...
func (i Item) DrawStatusLine(compactRendering bool, printedLineLength, realLineLength int) []string {
	if i.isRaw {
		return splitStatusLineToLines(i.CmdLine, printedLineLength, realLineLength)
	}
	secs := int64(i.realtimeBefore)
	nsecs := int64((i.realtimeBefore - float64(secs)) * 1e9)
	tm := time.Unix(secs, nsecs)
	const timeFormat = "2006-01-02 15:04:05"
	timeString := tm.Format(timeFormat)

	pwdTilde := strings.Replace(i.pwd, i.home, "~", 1)

	separator := "    "
	stLine := timeString + separator + i.host + ":" + pwdTilde + separator + i.CmdLine
	return splitStatusLineToLines(stLine, printedLineLength, realLineLength)
}

// GetEmptyStatusLine .
func GetEmptyStatusLine(printedLineLength, realLineLength int) []string {
	return splitStatusLineToLines("- no result selected -", printedLineLength, realLineLength)
}

// DrawItemColumns ...
func (i Item) DrawItemColumns(compactRendering bool, debug bool) ItemColumns {
	if i.isRaw {
		notAvailable := "n/a"
		return ItemColumns{
			Date:          notAvailable + " ",
			DateWithColor: notAvailable + " ",
			// dateWithColor:    highlightDate(notAvailable) + " ",
			Host:             "",
			HostWithColor:    "",
			PwdTilde:         notAvailable,
			CmdLine:          i.CmdLine,
			CmdLineWithColor: i.CmdLineWithColor,
			// score:             i.score,
			Key: i.Key,
		}
	}

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
		hitsStr := fmt.Sprintf("%.1f", i.Score)
		flags += " S" + hitsStr
		flagsWithColor += " S" + hitsStr
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
	return ItemColumns{
		Date:             date,
		DateWithColor:    dateWithColor,
		Host:             host,
		HostWithColor:    hostWithColor,
		PwdTilde:         pwdTilde,
		samePwd:          i.samePwd,
		Flags:            flags,
		FlagsWithColor:   flagsWithColor,
		CmdLine:          i.CmdLine,
		CmdLineWithColor: i.CmdLineWithColor,
		// score:             i.score,
		Key: i.Key,
	}
}

// ProduceLine ...
func (ic ItemColumns) ProduceLine(dateLength int, locationLength int, flagLength int, header bool, showDate bool) (string, int) {
	line := ""
	if showDate {
		date := ic.Date
		for len(date) < dateLength {
			line += " "
			date += " "
		}
		// TODO: use strings.Repeat
		line += ic.DateWithColor
	}
	// LOCATION
	locationWithColor := ic.HostWithColor
	pwdLength := locationLength - len(ic.Host)
	if ic.samePwd {
		locationWithColor += highlightPwd(leftCutPadString(ic.PwdTilde, pwdLength))
	} else {
		locationWithColor += leftCutPadString(ic.PwdTilde, pwdLength)
	}
	line += locationWithColor
	line += ic.FlagsWithColor
	flags := ic.Flags
	if flagLength < len(ic.Flags) {
		log.Printf("produceLine can't specify line w/ flags shorter than the actual size. - len(flags) %v, requested %v\n", len(ic.Flags), flagLength)
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
	line += spacer + ic.CmdLineWithColor

	length := dateLength + locationLength + flagLength + len(spacer) + len(ic.CmdLine)
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

// NewItemFromRecordForQuery creates new item from record based on given query
//		returns error if the query doesn't match the record
func NewItemFromRecordForQuery(record records.CliRecord, query Query, debug bool) (Item, error) {
	// Use numbers that won't add up to same score for any number of query words
	// query score weigth 1.51
	const hitScore = 1.517              // 1 * 1.51
	const properMatchScore = 0.501      // 0.33 * 1.51
	const hitScoreConsecutive = 0.00302 // 0.002 * 1.51

	// context score weigth 1
	// Host penalty
	var actualPwdScore = 0.9
	var sameGitRepoScore = 0.8
	var nonZeroExitCodeScorePenalty = 0.4
	var differentHostScorePenalty = 0.2

	reduceHostPenalty := false
	if reduceHostPenalty {
		actualPwdScore = 0.9
		sameGitRepoScore = 0.7
		nonZeroExitCodeScorePenalty = 0.4
		differentHostScorePenalty = 0.1
	}

	const timeScoreCoef = 1e-13
	// nonZeroExitCodeScorePenalty + differentHostScorePenalty

	score := 0.0
	anyHit := false
	cmd := record.CmdLine
	for _, term := range query.terms {
		c := strings.Count(record.CmdLine, term)
		if c > 0 {
			anyHit = true
			score += hitScore + hitScoreConsecutive*float64(c)
			if properMatch(cmd, term, " ") {
				score += properMatchScore
			}
			cmd = strings.ReplaceAll(cmd, term, highlightMatch(term))
		}
	}
	// DISPLAY > cmdline

	// cmd := "<" + strings.ReplaceAll(record.CmdLine, "\n", ";") + ">"
	cmdLine := strings.ReplaceAll(record.CmdLine, "\n", ";")
	cmdLineWithColor := strings.ReplaceAll(cmd, "\n", ";")

	// KEY for deduplication

	key := record.CmdLine
	// NOTE: since we import standard history we need a compatible key without metadata
	/*
		unlikelySeparator := "|||||"
		key := record.CmdLine + unlikelySeparator + record.Pwd + unlikelySeparator +
		record.GitOriginRemote + unlikelySeparator + record.Host
	*/
	if record.IsRaw {
		return Item{
			isRaw: true,

			CmdLine:          cmdLine,
			CmdLineWithColor: cmdLineWithColor,
			Score:            score,
			Key:              key,
		}, nil
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
	_ = anyHit
	// if score <= 0 && !anyHit {
	//	return Item{}, errors.New("no match for given record and query")
	// }
	score += record.RealtimeBefore * timeScoreCoef

	it := Item{
		realtimeBefore: record.RealtimeBefore,

		differentHost: differentHost,
		host:          record.Host,
		home:          record.Home,
		samePwd:       samePwd,
		pwd:           record.Pwd,

		sameGitRepo:      sameGitRepo,
		exitCode:         record.ExitCode,
		CmdLine:          cmdLine,
		CmdLineWithColor: cmdLineWithColor,
		Score:            score,
		Key:              key,
	}
	return it, nil
}

// GetHeader returns header columns
func GetHeader(compactRendering bool) ItemColumns {
	date := "TIME "
	host := "HOST:"
	dir := "DIRECTORY"
	if compactRendering {
		dir = "DIR"
	}
	flags := " FLAGS"
	cmdLine := "COMMAND-LINE"
	return ItemColumns{
		Date:             date,
		DateWithColor:    date,
		Host:             host,
		HostWithColor:    host,
		PwdTilde:         dir,
		samePwd:          false,
		Flags:            flags,
		FlagsWithColor:   flags,
		CmdLine:          cmdLine,
		CmdLineWithColor: cmdLine,
		// score:             i.score,
		Key: "_HEADERS_",
	}
}

// RawItem is item for raw mode
type RawItem struct {
	CmdLineWithColor string
	CmdLine          string

	Score float64

	Key string
	// cmdLineRaw string
}

// NewRawItemFromRecordForQuery creates new item from record based on given query
//		returns error if the query doesn't match the record
func NewRawItemFromRecordForQuery(record records.CliRecord, terms []string, debug bool) (RawItem, error) {
	const hitScore = 1.0
	const hitScoreConsecutive = 0.01
	const properMatchScore = 0.3

	const timeScoreCoef = 1e-13

	score := 0.0
	cmd := record.CmdLine
	for _, term := range terms {
		c := strings.Count(record.CmdLine, term)
		if c > 0 {
			score += hitScore + hitScoreConsecutive*float64(c)
			if properMatch(cmd, term, " ") {
				score += properMatchScore
			}
			cmd = strings.ReplaceAll(cmd, term, highlightMatch(term))
		}
	}
	score += record.RealtimeBefore * timeScoreCoef
	// KEY for deduplication
	key := record.CmdLine

	// DISPLAY > cmdline

	// cmd := "<" + strings.ReplaceAll(record.CmdLine, "\n", ";") + ">"
	cmdLine := strings.ReplaceAll(record.CmdLine, "\n", ";")
	cmdLineWithColor := strings.ReplaceAll(cmd, "\n", ";")

	it := RawItem{
		CmdLine:          cmdLine,
		CmdLineWithColor: cmdLineWithColor,
		Score:            score,
		Key:              key,
	}
	return it, nil
}
