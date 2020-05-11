package main

import (
	"bytes"
	"encoding/json"
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

	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	// g.SelFgColor = gocui.ColorGreen
	// g.SelBgColor = gocui.ColorGreen
	g.Highlight = true

	mess := msg.CliMsg{
		SessionID: *sessionID,
		PWD:       *pwd,
	}
	resp := SendCliMsg(mess, strconv.Itoa(config.Port))

	st := state{
		// lock sync.Mutex
		cliRecords:   resp.CliRecords,
		initialQuery: *query,
	}

	layout := manager{
		sessionID:       *sessionID,
		host:            *host,
		pwd:             *pwd,
		gitOriginRemote: records.NormalizeGitRemote(*gitOriginRemote),
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
	if err := g.SetKeybinding("", gocui.KeyCtrlN, gocui.ModNone, layout.Next); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, layout.Prev); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, layout.Prev); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, layout.SelectPaste); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, layout.SelectExecute); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlG, gocui.ModNone, layout.AbortPaste); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, layout.SwitchModes); err != nil {
		log.Panicln(err)
	}

	layout.UpdateData(*query)
	layout.UpdateRawData(*query)
	err = g.MainLoop()
	if err != nil && gocui.IsQuit(err) == false {
		log.Panicln(err)
	}
	return layout.s.output, layout.s.exitCode
}

type state struct {
	lock                sync.Mutex
	cliRecords          []records.CliRecord
	data                []item
	rawData             []rawItem
	highlightedItem     int
	displayedItemsCount int

	rawMode bool

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

func (m manager) AbortPaste(g *gocui.Gui, v *gocui.View) error {
	m.s.lock.Lock()
	defer m.s.lock.Unlock()
	if m.s.highlightedItem < len(m.s.data) {
		m.s.output = v.Buffer()
		m.s.exitCode = 0 // success
		return gocui.ErrQuit
	}
	return nil
}

type dedupRecord struct {
	dataIndex int
	score     float32
}

func (m manager) UpdateData(input string) {
	if debug {
		log.Println("EDIT start")
		log.Println("len(fullRecords) =", len(m.s.cliRecords))
		log.Println("len(data) =", len(m.s.data))
	}
	query := newQueryFromString(input, m.host, m.pwd, m.gitOriginRemote)
	var data []item
	itemSet := make(map[string]int)
	m.s.lock.Lock()
	defer m.s.lock.Unlock()
	for _, rec := range m.s.cliRecords {
		itm, err := newItemFromRecordForQuery(rec, query, m.config.Debug)
		if err != nil {
			// records didn't match the query
			// log.Println(" * continue (no match)", rec.Pwd)
			continue
		}
		if idx, ok := itemSet[itm.key]; ok {
			// duplicate found
			if data[idx].score >= itm.score {
				// skip duplicate item
				continue
			}
			// update duplicate item
			data[idx] = itm
			continue
		}
		// add new item
		itemSet[itm.key] = len(data)
		data = append(data, itm)
	}
	if debug {
		log.Println("len(tmpdata) =", len(data))
	}
	sort.SliceStable(data, func(p, q int) bool {
		return data[p].score > data[q].score
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
		log.Println("len(fullRecords) =", len(m.s.cliRecords))
		log.Println("len(data) =", len(m.s.data))
		log.Println("EDIT end")
	}
}

func (m manager) UpdateRawData(input string) {
	if debug {
		log.Println("EDIT start")
		log.Println("len(fullRecords) =", len(m.s.cliRecords))
		log.Println("len(data) =", len(m.s.data))
	}
	query := getRawTermsFromString(input)
	var data []rawItem
	itemSet := make(map[string]bool)
	m.s.lock.Lock()
	defer m.s.lock.Unlock()
	for _, rec := range m.s.cliRecords {
		itm, err := newRawItemFromRecordForQuery(rec, query, m.config.Debug)
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
	m.s.rawData = nil
	for _, itm := range data {
		if len(m.s.rawData) > 420 {
			break
		}
		m.s.rawData = append(m.s.rawData, itm)
	}
	m.s.highlightedItem = 0
	if debug {
		log.Println("len(fullRecords) =", len(m.s.cliRecords))
		log.Println("len(data) =", len(m.s.data))
		log.Println("EDIT end")
	}
}
func (m manager) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	gocui.DefaultEditor.Edit(v, key, ch, mod)
	if m.s.rawMode {
		m.UpdateRawData(v.Buffer())
		return
	}
	m.UpdateData(v.Buffer())
}

func (m manager) Next(g *gocui.Gui, v *gocui.View) error {
	m.s.lock.Lock()
	defer m.s.lock.Unlock()
	if m.s.highlightedItem < m.s.displayedItemsCount-1 {
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

func (m manager) SwitchModes(g *gocui.Gui, v *gocui.View) error {
	m.s.lock.Lock()
	m.s.rawMode = !m.s.rawMode
	m.s.lock.Unlock()

	if m.s.rawMode {
		m.UpdateRawData(v.Buffer())
		return nil
	}
	m.UpdateData(v.Buffer())
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
	if m.s.rawMode {
		v.Title = " RESH CLI - NON-CONTEXTUAL \"RAW\" MODE - (CTRL+R to switch BACK) "
	} else {
		v.Title = " RESH CLI - CONTEXTUAL MODE - (CTRL+R to switch to RAW MODE) "
	}

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

	if m.s.rawMode {
		return m.rawMode(g, v)
	}
	return m.normalMode(g, v)
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func getHeader(compactRendering bool) itemColumns {
	date := "TIME "
	host := "HOST:"
	dir := "DIRECTORY"
	if compactRendering {
		dir = "DIR"
	}
	flags := " FLAGS"
	cmdLine := "COMMAND-LINE"
	return itemColumns{
		date:             date,
		dateWithColor:    date,
		host:             host,
		hostWithColor:    host,
		pwdTilde:         dir,
		samePwd:          false,
		flags:            flags,
		flagsWithColor:   flags,
		cmdLine:          cmdLine,
		cmdLineWithColor: cmdLine,
		// score:             i.score,
		key: "_HEADERS_",
	}
}

const smallTerminalTresholdWidth = 110

func (m manager) normalMode(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()

	compactRenderingMode := false
	if maxX < smallTerminalTresholdWidth {
		compactRenderingMode = true
	}

	data := []itemColumns{}

	header := getHeader(compactRenderingMode)
	longestDateLen := len(header.date)
	longestLocationLen := len(header.host) + len(header.pwdTilde)
	longestFlagsLen := 2
	maxPossibleMainViewHeight := maxY - 3 - 1 - 1 - 1 // - top box - header - status - help
	for i, itm := range m.s.data {
		if i == maxY {
			break
		}
		ic := itm.drawItemColumns(compactRenderingMode)
		data = append(data, ic)
		if i > maxPossibleMainViewHeight {
			// do not stretch columns because of results that will end up outside of the page
			continue
		}
		if len(ic.date) > longestDateLen {
			longestDateLen = len(ic.date)
		}
		if len(ic.host)+len(ic.pwdTilde) > longestLocationLen {
			longestLocationLen = len(ic.host) + len(ic.pwdTilde)
		}
		if len(ic.flags) > longestFlagsLen {
			longestFlagsLen = len(ic.flags)
		}
	}
	maxLocationLen := maxX/7 + 8
	if longestLocationLen > maxLocationLen {
		longestLocationLen = maxLocationLen
	}

	if m.s.highlightedItem >= len(m.s.data) {
		m.s.highlightedItem = len(m.s.data) - 1
	}
	// status line
	topBoxHeight := 3 // size of the query box up top
	topBoxHeight++    // headers
	realLineLength := maxX - 2
	printedLineLength := maxX - 4
	statusLine := m.s.data[m.s.highlightedItem].drawStatusLine(compactRenderingMode, printedLineLength, realLineLength)
	var statusLineHeight int = len(statusLine) + 1 // help line

	helpLineHeight := 1
	const helpLine = "HELP: type to search, UP/DOWN to select, RIGHT to edit, ENTER to execute, CTRL+G to abort, CTRL+C/D to quit; " +
		"TIP: when resh-cli is launched command line is used as initial search query"

	mainViewHeight := maxY - topBoxHeight - statusLineHeight - helpLineHeight
	m.s.displayedItemsCount = mainViewHeight

	// header
	// header := getHeader()
	dispStr, _ := header.produceLine(longestDateLen, longestLocationLen, longestFlagsLen, true, true)
	dispStr = doHighlightHeader(dispStr, maxX*2)
	v.WriteString(dispStr + "\n")

	var index int
	for index < len(data) {
		itm := data[index]
		if index == mainViewHeight {
			// page is full
			break
		}

		displayStr, _ := itm.produceLine(longestDateLen, longestLocationLen, longestFlagsLen, false, true)
		if m.s.highlightedItem == index {
			// maxX * 2 because there are escape sequences that make it hard to tell the real string lenght
			displayStr = doHighlightString(displayStr, maxX*3)
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
		index++
	}
	// push the status line to the bottom of the page
	for index < mainViewHeight {
		v.WriteString("\n")
		index++
	}
	for _, line := range statusLine {
		v.WriteString(line)
	}
	v.WriteString(helpLine)
	if debug {
		log.Println("len(data) =", len(m.s.data))
		log.Println("highlightedItem =", m.s.highlightedItem)
	}
	return nil
}

func (m manager) rawMode(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	topBoxSize := 3
	m.s.displayedItemsCount = maxY - topBoxSize

	for i, itm := range m.s.rawData {
		if i == maxY {
			if debug {
				log.Println(maxY)
			}
			break
		}
		displayStr := itm.cmdLineWithColor
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

// SendCliMsg to daemon
func SendCliMsg(m msg.CliMsg, port string) msg.CliResponse {
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
	response := msg.CliResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal("unmarshal resp error: ", err)
	}
	return response
}
