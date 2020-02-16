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
	}

	sessionID := flag.String("sessionID", "", "resh generated session id")
	pwd := flag.String("pwd", "", "present working directory")
	query := flag.String("query", "", "search query")
	flag.Parse()

	if *sessionID == "" {
		fmt.Println("Error: you need to specify sessionId")
	}
	if *pwd == "" {
		fmt.Println("Error: you need to specify PWD")
	}

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
		sessionID: *sessionID,
		pwd:       *pwd,
		config:    config,
		s:         &st,
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
	blueBold := "\033[34;1m"
	redBold := "\033[31;1m"
	repace := []string{invert, end, blueBold, redBold}
	if strings.Contains(str, prefix) == false {
		return str
	}
	for _, escSeq := range repace {
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

func highlightMatchAlternative(str string) string {
	// template "\033[3%d;%dm"
	blueBold := "\033[34;1m"
	end := "\033[0m"
	return blueBold + cleanHighlight(str) + end
}

func highlightMatch(str string) string {
	// template "\033[3%d;%dm"
	redBold := "\033[31;1m"
	end := "\033[0m"
	return redBold + cleanHighlight(str) + end
}

func toString(record records.EnrichedRecord, lineLength int) string {
	dirColWidth := 24 // make this dynamic somehow
	return leftCutPadString(strings.Replace(record.Pwd, record.Home, "~", 1), dirColWidth) + "   " +
		rightCutPadString(strings.ReplaceAll(record.CmdLine, "\n", "; "), lineLength-dirColWidth-3) + "\n"
}

type query struct {
	terms []string
	pwd   string
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

func newQueryFromString(queryInput string, pwd string) query {
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
	return query{terms: terms, pwd: pwd}
}

type item struct {
	// record         records.EnrichedRecord
	display        string
	displayNoColor string
	cmdLine        string
	pwd            string
	pwdTilde       string
	hits           float64
}

func (i item) less(i2 item) bool {
	// reversed order
	return i.hits > i2.hits
}

// used for deduplication
func (i item) key() string {
	unlikelySeparator := "|||||"
	return i.cmdLine + unlikelySeparator + i.pwd
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
	// TODO: use color to highlight matches
	const hitScore = 1.0
	const hitScoreConsecutive = 0.1
	const properMatchScore = 0.3
	const actualPwdScore = 0.9
	const actualPwdScoreExtra = 0.2

	hits := 0.0
	if record.ExitCode != 0 {
		hits--
	}
	cmd := record.CmdLine
	pwdTilde := strings.Replace(record.Pwd, record.Home, "~", 1)
	pwdDisp := leftCutPadString(pwdTilde, 25)
	pwdRawDisp := leftCutPadString(record.Pwd, 25)
	var useRawPwd bool
	var dirHit bool
	for _, term := range query.terms {
		termHit := false
		if strings.Contains(record.CmdLine, term) {
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
		if strings.Contains(pwdTilde, term) {
			if termHit == false {
				hits += hitScore
			} else {
				hits += hitScoreConsecutive
			}
			termHit = true
			if properMatch(pwdTilde, term, "/") {
				hits += properMatchScore
			}
			pwdDisp = strings.ReplaceAll(pwdDisp, term, highlightMatch(term))
			pwdRawDisp = strings.ReplaceAll(pwdRawDisp, term, highlightMatch(term))
			dirHit = true
		} else if strings.Contains(record.Pwd, term) {
			if termHit == false {
				hits += hitScore
			} else {
				hits += hitScoreConsecutive
			}
			termHit = true
			if properMatch(pwdTilde, term, "/") {
				hits += properMatchScore
			}
			pwdRawDisp = strings.ReplaceAll(pwdRawDisp, term, highlightMatch(term))
			dirHit = true
			useRawPwd = true
		}
		// if strings.Contains(record.GitOriginRemote, term) {
		// 	hits++
		// }
	}
	// actual pwd matches
	// only use if there was no directory match on any of the terms
	// N terms can only produce:
	//		-> N matches against the command
	//		-> N matches against the directory
	//		-> 1 extra match for the actual directory match
	if record.Pwd == query.pwd {
		if dirHit {
			hits += actualPwdScoreExtra
		} else {
			hits += actualPwdScore
		}
		pwdDisp = highlightMatchAlternative(pwdDisp)
		// pwdRawDisp = highlightMatchAlternative(pwdRawDisp)
		useRawPwd = false
	}
	if hits <= 0 {
		return item{}, errors.New("no match for given record and query")
	}
	display := ""
	// pwd := leftCutPadString("<"+pwdTilde+">", 20)
	if useRawPwd {
		display += pwdRawDisp
	} else {
		display += pwdDisp
	}
	if debug {
		hitsStr := fmt.Sprintf("%.1f", hits)
		hitsDisp := "  " + hitsStr + "  "
		display += hitsDisp
	} else {
		display += "  "
	}
	// cmd := "<" + strings.ReplaceAll(record.CmdLine, "\n", ";") + ">"
	cmd = strings.ReplaceAll(cmd, "\n", ";")
	display += cmd
	// itDummy := item{
	// 	cmdLine: record.CmdLine,
	// 	pwd:     record.Pwd,
	// }
	// + "   #K:<" + itDummy.key() + ">"

	it := item{
		display:        display,
		displayNoColor: display,
		cmdLine:        record.CmdLine,
		pwd:            record.Pwd,
		pwdTilde:       pwdTilde,
		hits:           hits,
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
	sessionID string
	pwd       string
	config    cfg.Config

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
	query := newQueryFromString(input, m.pwd)
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
		if itemSet[itm.key()] {
			// log.Println(" * continue (already present)", itm.key(), itm.pwd)
			continue
		}
		itemSet[itm.key()] = true
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

	for i, itm := range m.s.data {
		if i == maxY {
			log.Println(maxY)
			break
		}
		displayStr := itm.display
		if m.s.highlightedItem == i {
			// use actual min requried length instead of 420 constant
			displayStr = doHighlightString(displayStr, 420)
			log.Println("### HightlightedItem string :", displayStr)
		} else {
			log.Println(displayStr)
		}
		if strings.Contains(displayStr, "\n") {
			log.Println("display string contained \\n")
			displayStr = strings.ReplaceAll(displayStr, "\n", "#")
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
