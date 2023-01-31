package searchapp

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/curusarn/resh/internal/recordint"
	"golang.org/x/exp/utf8string"
)

const itemLocationLength = 30
const dots = "…"

// Item holds item info for normal mode
type Item struct {
	isRaw bool

	time float64

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
	differentHost bool
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
	secs := int64(i.time)
	nsecs := int64((i.time - float64(secs)) * 1e9)
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
			PwdTilde:         notAvailable,
			CmdLine:          i.CmdLine,
			CmdLineWithColor: i.CmdLineWithColor,
			// score:             i.score,
			Key: i.Key,
		}
	}

	// DISPLAY
	// DISPLAY > date
	secs := int64(i.time)
	nsecs := int64((i.time - float64(secs)) * 1e9)
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
	if i.differentHost {
		host += i.host
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
		PwdTilde:         pwdTilde,
		samePwd:          i.samePwd,
		differentHost:    i.differentHost,
		Flags:            flags,
		FlagsWithColor:   flagsWithColor,
		CmdLine:          i.CmdLine,
		CmdLineWithColor: i.CmdLineWithColor,
		// score:             i.score,
		Key: i.Key,
	}
}

func minInt(values ...int) int {
	min := math.MaxInt32
	for _, val := range values {
		if val < min {
			min = val
		}
	}
	return min
}

func produceLocation(length int, host string, pwdTilde string, differentHost bool, samePwd bool, debug bool) string {
	hostLen := len(host)
	if hostLen <= 0 {
		pwdWithColor := leftCutPadString(pwdTilde, length)
		if samePwd {
			pwdWithColor = highlightPwd(pwdWithColor)
		}
		return pwdWithColor
	}
	colonLen := 1
	pwdLen := len(pwdTilde)
	totalLen := hostLen + colonLen + pwdLen

	newHostLen := hostLen
	// only shrink if the location does not fit
	if totalLen > length {
		// how much we need to shrink/crop the location
		shrinkFactor := float64(length) / float64(totalLen)

		shrinkedHostLen := int(math.Ceil(float64(hostLen) * shrinkFactor))
		halfLocationLen := length/2 - colonLen

		newHostLen = minInt(hostLen, shrinkedHostLen, halfLocationLen)
	}
	// pwd length is the rest of the length
	newPwdLen := length - colonLen - newHostLen

	// adjust pwd length
	if newPwdLen > pwdLen {
		diff := newPwdLen - pwdLen
		newHostLen += diff
		newPwdLen -= diff
	}

	hostWithColor := rightCutLeftPadString(host, newHostLen)
	if differentHost {
		hostWithColor = highlightHost(hostWithColor)
	}
	pwdWithColor := leftCutPadString(pwdTilde, newPwdLen)
	if samePwd {
		pwdWithColor = highlightPwd(pwdWithColor)
	}
	return hostWithColor + ":" + pwdWithColor
}

// ProduceLine ...
func (ic ItemColumns) ProduceLine(dateLength int, locationLength int, flagsLength int, header bool, showDate bool, debug bool) (string, int, error) {
	var err error
	line := ""
	if showDate {
		line += strings.Repeat(" ", dateLength-len(ic.Date)) + ic.DateWithColor
	}
	// LOCATION
	locationWithColor := produceLocation(locationLength, ic.Host, ic.PwdTilde, ic.differentHost, ic.samePwd, debug)
	line += locationWithColor

	// FLAGS
	line += ic.FlagsWithColor
	if flagsLength >= len(ic.Flags) {
		line += strings.Repeat(" ", flagsLength-len(ic.Flags))
	} else {
		err = fmt.Errorf("actual flags are longer than dedicated flag space. actual: %v, space: %v", len(ic.Flags), flagsLength)
	}
	spacer := "  "
	if flagsLength > 5 || header {
		// use shorter spacer
		// 		because there is likely a long flag like E130 in the view
		spacer = " "
	}
	line += spacer + ic.CmdLineWithColor

	length := dateLength + locationLength + flagsLength + len(spacer) + len(ic.CmdLine)
	return line, length, err
}

func rightCutLeftPadString(str string, newLen int) string {
	if newLen <= 0 {
		return ""
	}
	utf8Str := utf8string.NewString(str)
	strLen := utf8Str.RuneCount()
	if newLen > strLen {
		return strings.Repeat(" ", newLen-strLen) + str
	} else if newLen < strLen {
		return utf8Str.Slice(0, newLen-1) + dots
	}
	return str
}

func leftCutPadString(str string, newLen int) string {
	if newLen <= 0 {
		return ""
	}
	utf8Str := utf8string.NewString(str)
	strLen := utf8Str.RuneCount()
	if newLen > strLen {
		return strings.Repeat(" ", newLen-strLen) + str
	} else if newLen < strLen {
		return dots + utf8string.NewString(str).Slice(strLen-newLen+1, strLen)
	}
	return str
}

func rightCutPadString(str string, newLen int) string {
	if newLen <= 0 {
		return ""
	}
	utf8Str := utf8string.NewString(str)
	strLen := utf8Str.RuneCount()
	if newLen > strLen {
		return str + strings.Repeat(" ", newLen-strLen)
	} else if newLen < strLen {
		return utf8Str.Slice(0, newLen-1) + dots
	}
	return str
}

// proper match for path is when whole directory is matched
// proper match for command is when term matches word delimited by whitespace
func properMatch(str, term, padChar string) bool {
	return strings.Contains(padChar+str+padChar, padChar+term+padChar)
}

// NewItemFromRecordForQuery creates new item from record based on given query
//
//	returns error if the query doesn't match the record
func NewItemFromRecordForQuery(record recordint.SearchApp, query Query, debug bool) (Item, error) {
	// Use numbers that won't add up to same score for any number of query words
	// query score weight 1.51
	const hitScore = 1.517              // 1 * 1.51
	const properMatchScore = 0.501      // 0.33 * 1.51
	const hitScoreConsecutive = 0.00302 // 0.002 * 1.51

	// context score weight 1
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
	score += record.Time * timeScoreCoef

	it := Item{
		time: record.Time,

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
	host := "HOST"
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
//
//	returns error if the query doesn't match the record
func NewRawItemFromRecordForQuery(record recordint.SearchApp, terms []string, debug bool) (RawItem, error) {
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
	score += record.Time * timeScoreCoef
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
