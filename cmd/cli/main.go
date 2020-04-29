package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/awesome-gocui/gocui"
	"github.com/curusarn/resh/pkg/cfg"
	"github.com/curusarn/resh/pkg/msg"
	"github.com/curusarn/resh/pkg/records"

	"os/user"
	"path/filepath"
	"strconv"
)

// version from git set during build
var version string

// commit from git set during build
var commit string

// special constant recognized by RESH wrappers
const exitCodeExecute = 111

var debug bool

func main() {
	output, exitCode := runReshCli()
	fmt.Print(output)
	os.Exit(exitCode)
}

func runReshCli() (string, int) {
	usr, _ := user.Current()
	dir := usr.HomeDir
	configPath := filepath.Join(dir, "/.config/resh.toml")
	logPath := filepath.Join(dir, ".resh/cli.log")

	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer f.Close()

	log.SetOutput(f)

	var config cfg.Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Fatal("Error reading config:", err)
	}
	if config.Debug {
		debug = true
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
		log.Println("DEBUG is ON")
	}

	sessionID := flag.String("sessionID", "", "resh generated session id")
	host := flag.String("host", "", "host")
	pwd := flag.String("pwd", "", "present working directory")
	gitOriginRemote := flag.String("gitOriginRemote", "DEFAULT", "git origin remote")
	query := flag.String("query", "", "search query")
	flag.Parse()

	if *sessionID == "" {
		log.Println("Error: you need to specify sessionId")
	}
	if *host == "" {
		log.Println("Error: you need to specify HOST")
	}
	if *pwd == "" {
		log.Println("Error: you need to specify PWD")
	}
	if *gitOriginRemote == "DEFAULT" {
		log.Println("Error: you need to specify gitOriginRemote")
	}

	log.Printf("gitRemoteOrigin: %s\n", *gitOriginRemote)
	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	// g.SelBgColor = gocui.ColorGreen
	g.Highlight = true

	mess := msg.DumpMsg{
		SessionID: *sessionID,
		PWD:       *pwd,
	}
	resp := SendDumpMsg(mess, strconv.Itoa(config.Port))

	st := state{
		// lock sync.Mutex
		fullRecords:  resp.FullRecords,
		initialQuery: *query,
	}

	layout := manager{
		sessionID:       *sessionID,
		host:            *host,
		pwd:             *pwd,
		gitOriginRemote: *gitOriginRemote,
		config:          config,
		s:               &st,
	}
	g.SetManager(layout)

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, layout.Next); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, layout.Next); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, layout.Prev); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlG, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, layout.SelectExecute); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, layout.SelectPaste); err != nil {
		log.Panicln(err)
	}

	layout.UpdateData(*query)
	err = g.MainLoop()
	if err != nil && gocui.IsQuit(err) == false {
		log.Panicln(err)
	}
	return layout.s.output, layout.s.exitCode
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

func cleanHighlight(str string) string {
	prefix := "\033["

	invert := "\033[32;7;1m"
	end := "\033[0m"
	replace := []string{invert, end}
	for i := 30; i < 48; i++ {
		base := prefix + strconv.Itoa(i)
		normal := base + "m"
		bold := base + ";1m"
		replace = append(replace, normal, bold)
	}
	if strings.Contains(str, prefix) == false {
		return str
	}
	for _, escSeq := range replace {
		str = strings.ReplaceAll(str, escSeq, "")
	}
	return str
}

func highlightSelected(str string) string {
	// template "\033[3%d;%dm"
	invert := "\033[32;7;1m"
	end := "\033[0m"
	return invert + cleanHighlight(str) + end
}

func highlightHost(str string) string {
	// template "\033[3%d;%dm"
	redNormal := "\033[31m"
	end := "\033[0m"
	return redNormal + cleanHighlight(str) + end
}

func highlightPwd(str string) string {
	// template "\033[3%d;%dm"
	blueBold := "\033[34;1m"
	end := "\033[0m"
	return blueBold + cleanHighlight(str) + end
}

func highlightMatch(str string) string {
	// template "\033[3%d;%dm"
	magentaBold := "\033[35;1m"
	end := "\033[0m"
	return magentaBold + cleanHighlight(str) + end
}

func highlightWarn(str string) string {
	// template "\033[3%d;%dm"
	// orangeBold := "\033[33;1m"
	redBold := "\033[31;1m"
	end := "\033[0m"
	return redBold + cleanHighlight(str) + end
}

func highlightGit(str string) string {
	// template "\033[3%d;%dm"
	greenBold := "\033[32;1m"
	end := "\033[0m"
	return greenBold + cleanHighlight(str) + end
}

func toString(record records.EnrichedRecord, lineLength int) string {
	dirColWidth := 24 // make this dynamic somehow
	return leftCutPadString(strings.Replace(record.Pwd, record.Home, "~", 1), dirColWidth) + "   " +
		rightCutPadString(strings.ReplaceAll(record.CmdLine, "\n", "; "), lineLength-dirColWidth-3) + "\n"
}

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

	hits float64

	key string
	// cmdLineRaw string
}

func (i item) less(i2 item) bool {
	// reversed order
	return i.hits > i2.hits
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

// func (i item) equals(i2 item) bool {
// 	return i.cmdLine == i2.cmdLine && i.pwd == i2.pwd
// }

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
func newItemFromRecordForQuery(record records.EnrichedRecord, query query, debug bool) (item, error) {
	const hitScore = 1.0
	const hitScoreConsecutive = 0.1
	const properMatchScore = 0.3
	const actualPwdScore = 0.9
	const nonZeroExitCodeScorePenalty = 0.5
	const sameGitRepoScore = 0.7
	// const sameGitRepoScoreExtra = 0.0
	const differentHostScorePenalty = 0.2

	// nonZeroExitCodeScorePenalty + differentHostScorePenalty

	hits := 0.0
	anyHit := false
	cmd := record.CmdLine
	for _, term := range query.terms {
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
		hits += actualPwdScore
	} else if sameGitRepo {
		anyHit = true
		hits += sameGitRepoScore
	}

	differentHost := false
	if record.Host != query.host {
		differentHost = true
		hits -= differentHostScorePenalty
	}
	errorExitStatus := false
	if record.ExitCode != 0 {
		errorExitStatus = true
		hits -= nonZeroExitCodeScorePenalty
	}
	if hits <= 0 && !anyHit {
		return item{}, errors.New("no match for given record and query")
	}

	// KEY for deduplication

	unlikelySeparator := "|||||"
	key := record.CmdLine + unlikelySeparator + record.Pwd +
		unlikelySeparator + strconv.Itoa(record.ExitCode) + unlikelySeparator +
		record.GitOriginRemote + unlikelySeparator + record.Host

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
		hitsStr := fmt.Sprintf("%.1f", hits)
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
		hits:              hits,
		key:               key,
	}
	return it, nil
}

func doHighlightString(str string, minLength int) string {
	if len(str) < minLength {
		str = str + strings.Repeat(" ", minLength-len(str))
	}
	return highlightSelected(str)
}

type state struct {
	lock            sync.Mutex
	fullRecords     []records.EnrichedRecord
	data            []item
	highlightedItem int

	initialQuery string

	output   string
	exitCode int
}

type manager struct {
	sessionID       string
	host            string
	pwd             string
	gitOriginRemote string
	config          cfg.Config

	s *state
}

func (m manager) SelectExecute(g *gocui.Gui, v *gocui.View) error {
	m.s.lock.Lock()
	defer m.s.lock.Unlock()
	if m.s.highlightedItem < len(m.s.data) {
		m.s.output = m.s.data[m.s.highlightedItem].cmdLine
		m.s.exitCode = exitCodeExecute
		return gocui.ErrQuit
	}
	return nil
}

func (m manager) SelectPaste(g *gocui.Gui, v *gocui.View) error {
	m.s.lock.Lock()
	defer m.s.lock.Unlock()
	if m.s.highlightedItem < len(m.s.data) {
		m.s.output = m.s.data[m.s.highlightedItem].cmdLine
		m.s.exitCode = 0 // success
		return gocui.ErrQuit
	}
	return nil
}

func (m manager) UpdateData(input string) {
	if debug {
		log.Println("EDIT start")
		log.Println("len(fullRecords) =", len(m.s.fullRecords))
		log.Println("len(data) =", len(m.s.data))
	}
	query := newQueryFromString(input, m.host, m.pwd, m.gitOriginRemote)
	var data []item
	itemSet := make(map[string]bool)
	m.s.lock.Lock()
	defer m.s.lock.Unlock()
	for _, rec := range m.s.fullRecords {
		itm, err := newItemFromRecordForQuery(rec, query, m.config.Debug)
		if err != nil {
			// records didn't match the query
			// log.Println(" * continue (no match)", rec.Pwd)
			continue
		}
		if itemSet[itm.key] {
			// log.Println(" * continue (already present)", itm.key(), itm.pwd)
			continue
		}
		itemSet[itm.key] = true
		data = append(data, itm)
		// log.Println("DATA =", itm.display)
	}
	if debug {
		log.Println("len(tmpdata) =", len(data))
	}
	sort.SliceStable(data, func(p, q int) bool {
		return data[p].hits > data[q].hits
	})
	m.s.data = nil
	for _, itm := range data {
		if len(m.s.data) > 420 {
			break
		}
		m.s.data = append(m.s.data, itm)
	}
	m.s.highlightedItem = 0
	if debug {
		log.Println("len(fullRecords) =", len(m.s.fullRecords))
		log.Println("len(data) =", len(m.s.data))
		log.Println("EDIT end")
	}
}

func (m manager) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	gocui.DefaultEditor.Edit(v, key, ch, mod)
	m.UpdateData(v.Buffer())
}

func (m manager) Next(g *gocui.Gui, v *gocui.View) error {
	_, y := g.Size()
	m.s.lock.Lock()
	defer m.s.lock.Unlock()
	if m.s.highlightedItem < y {
		m.s.highlightedItem++
	}
	return nil
}

func (m manager) Prev(g *gocui.Gui, v *gocui.View) error {
	m.s.lock.Lock()
	defer m.s.lock.Unlock()
	if m.s.highlightedItem > 0 {
		m.s.highlightedItem--
	}
	return nil
}

func (m manager) Layout(g *gocui.Gui) error {
	var b byte
	maxX, maxY := g.Size()

	v, err := g.SetView("input", 0, 0, maxX-1, 2, b)
	if err != nil && gocui.IsUnknownView(err) == false {
		log.Panicln(err.Error())
	}

	v.Editable = true
	v.Editor = m
	v.Title = "resh cli"

	g.SetCurrentView("input")

	m.s.lock.Lock()
	defer m.s.lock.Unlock()
	if len(m.s.initialQuery) > 0 {
		v.WriteString(m.s.initialQuery)
		v.SetCursor(len(m.s.initialQuery), 0)
		m.s.initialQuery = ""
	}

	v, err = g.SetView("body", 0, 2, maxX-1, maxY, b)
	if err != nil && gocui.IsUnknownView(err) == false {
		log.Panicln(err.Error())
	}
	v.Frame = false
	v.Autoscroll = false
	v.Clear()
	v.Rewind()

	longestFlagsLen := 2 // at least 2
	for i, itm := range m.s.data {
		if i == maxY {
			break
		}
		if len(itm.flags) > longestFlagsLen {
			longestFlagsLen = len(itm.flags)
		}
	}

	for i, itm := range m.s.data {
		if i == maxY {
			if debug {
				log.Println(maxY)
			}
			break
		}
		displayStr, _ := itm.produceLine(longestFlagsLen)
		if m.s.highlightedItem == i {
			// use actual min requried length instead of 420 constant
			displayStr = doHighlightString(displayStr, maxX*2)
			if debug {
				log.Println("### HightlightedItem string :", displayStr)
			}
		} else if debug {
			log.Println(displayStr)
		}
		if strings.Contains(displayStr, "\n") {
			log.Println("display string contained \\n")
			displayStr = strings.ReplaceAll(displayStr, "\n", "#")
			if debug {
				log.Println("display string contained \\n")
			}
		}
		v.WriteString(displayStr + "\n")
		// if m.s.highlightedItem == i {
		// 	v.SetHighlight(m.s.highlightedItem, true)
		// }
	}
	if debug {
		log.Println("len(data) =", len(m.s.data))
		log.Println("highlightedItem =", m.s.highlightedItem)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// SendDumpMsg to daemon
func SendDumpMsg(m msg.DumpMsg, port string) msg.DumpResponse {
	recJSON, err := json.Marshal(m)
	if err != nil {
		log.Fatal("send err 1", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:"+port+"/dump",
		bytes.NewBuffer(recJSON))
	if err != nil {
		log.Fatal("send err 2", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("resh-daemon is not running :(")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("read response error")
	}
	// log.Println(string(body))
	response := msg.DumpResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal("unmarshal resp error: ", err)
	}
	return response
}
