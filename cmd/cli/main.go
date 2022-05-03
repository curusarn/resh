package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/curusarn/resh/internal/cfg"
	"github.com/curusarn/resh/internal/logger"
	"github.com/curusarn/resh/internal/msg"
	"github.com/curusarn/resh/internal/output"
	"github.com/curusarn/resh/internal/records"
	"github.com/curusarn/resh/internal/searchapp"
	"go.uber.org/zap"

	"strconv"
)

// info passed during build
var version string
var commit string
var developement bool

// special constant recognized by RESH wrappers
const exitCodeExecute = 111

func main() {
	config, errCfg := cfg.New()
	logger, _ := logger.New("search-app", config.LogLevel, developement)
	defer logger.Sync() // flushes buffer, if any
	if errCfg != nil {
		logger.Error("Error while getting configuration", zap.Error(errCfg))
	}
	out := output.New(logger, "resh-search-app ERROR")

	output, exitCode := runReshCli(out, config)
	fmt.Print(output)
	os.Exit(exitCode)
}

func runReshCli(out *output.Output, config cfg.Config) (string, int) {
	sessionID := flag.String("sessionID", "", "resh generated session id")
	host := flag.String("host", "", "host")
	pwd := flag.String("pwd", "", "present working directory")
	gitOriginRemote := flag.String("gitOriginRemote", "DEFAULT", "git origin remote")
	query := flag.String("query", "", "search query")
	testHistory := flag.String("test-history", "", "load history from a file instead from the daemon (for testing purposes only!)")
	testHistoryLines := flag.Int("test-lines", 0, "the number of lines to load from a file passed with --test-history (for testing purposes only!)")
	flag.Parse()

	errMsg := "Failed to get necessary command-line arguments"
	if *sessionID == "" {
		out.Fatal(errMsg, errors.New("missing option --sessionId"))
	}
	if *host == "" {
		out.Fatal(errMsg, errors.New("missing option --host"))
	}
	if *pwd == "" {
		out.Fatal(errMsg, errors.New("missing option --pwd"))
	}
	if *gitOriginRemote == "DEFAULT" {
		out.Fatal(errMsg, errors.New("missing option --gitOriginRemote"))
	}

	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		out.Fatal("Failed to launch TUI", err)
	}
	defer g.Close()

	g.Cursor = true
	// g.SelFgColor = gocui.ColorGreen
	// g.SelBgColor = gocui.ColorGreen
	g.Highlight = true

	var resp msg.CliResponse
	if *testHistory == "" {
		mess := msg.CliMsg{
			SessionID: *sessionID,
			PWD:       *pwd,
		}
		resp = SendCliMsg(out, mess, strconv.Itoa(config.Port))
	} else {
		resp = searchapp.LoadHistoryFromFile(out.Logger.Sugar(), *testHistory, *testHistoryLines)
	}

	st := state{
		// lock sync.Mutex
		cliRecords:   resp.CliRecords,
		initialQuery: *query,
	}

	layout := manager{
		out:             out,
		config:          config,
		sessionID:       *sessionID,
		host:            *host,
		pwd:             *pwd,
		gitOriginRemote: records.NormalizeGitRemote(*gitOriginRemote),
		s:               &st,
	}
	g.SetManager(layout)

	errMsg = "Failed to set keybindings"
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, layout.Next); err != nil {
		out.Fatal(errMsg, err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, layout.Next); err != nil {
		out.Fatal(errMsg, err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlN, gocui.ModNone, layout.Next); err != nil {
		out.Fatal(errMsg, err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, layout.Prev); err != nil {
		out.Fatal(errMsg, err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, layout.Prev); err != nil {
		out.Fatal(errMsg, err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, layout.SelectPaste); err != nil {
		out.Fatal(errMsg, err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, layout.SelectExecute); err != nil {
		out.Fatal(errMsg, err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlG, gocui.ModNone, layout.AbortPaste); err != nil {
		out.Fatal(errMsg, err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		out.Fatal(errMsg, err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, quit); err != nil {
		out.Fatal(errMsg, err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, layout.SwitchModes); err != nil {
		out.Fatal(errMsg, err)
	}

	layout.UpdateData(*query)
	layout.UpdateRawData(*query)
	err = g.MainLoop()
	if err != nil && !errors.Is(err, gocui.ErrQuit) {
		out.Fatal("Main application loop finished with error", err)
	}
	return layout.s.output, layout.s.exitCode
}

type state struct {
	lock                sync.Mutex
	cliRecords          []records.CliRecord
	data                []searchapp.Item
	rawData             []searchapp.RawItem
	highlightedItem     int
	displayedItemsCount int

	rawMode bool

	initialQuery string

	output   string
	exitCode int
}

type manager struct {
	out    *output.Output
	config cfg.Config

	sessionID       string
	host            string
	pwd             string
	gitOriginRemote string

	s *state
}

func (m manager) SelectExecute(g *gocui.Gui, v *gocui.View) error {
	m.s.lock.Lock()
	defer m.s.lock.Unlock()
	if m.s.rawMode {
		if m.s.highlightedItem < len(m.s.rawData) {
			m.s.output = m.s.rawData[m.s.highlightedItem].CmdLine
			m.s.exitCode = exitCodeExecute
			return gocui.ErrQuit
		}
	} else {
		if m.s.highlightedItem < len(m.s.data) {
			m.s.output = m.s.data[m.s.highlightedItem].CmdLine
			m.s.exitCode = exitCodeExecute
			return gocui.ErrQuit
		}
	}
	return nil
}

func (m manager) SelectPaste(g *gocui.Gui, v *gocui.View) error {
	m.s.lock.Lock()
	defer m.s.lock.Unlock()
	if m.s.rawMode {
		if m.s.highlightedItem < len(m.s.rawData) {
			m.s.output = m.s.rawData[m.s.highlightedItem].CmdLine
			m.s.exitCode = 0 // success
			return gocui.ErrQuit
		}
	} else {
		if m.s.highlightedItem < len(m.s.data) {
			m.s.output = m.s.data[m.s.highlightedItem].CmdLine
			m.s.exitCode = 0 // success
			return gocui.ErrQuit
		}
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
	sugar := m.out.Logger.Sugar()
	sugar.Debugw("Starting data update ...",
		"recordCount", len(m.s.cliRecords),
		"itemCount", len(m.s.data),
	)
	query := searchapp.NewQueryFromString(input, m.host, m.pwd, m.gitOriginRemote, m.config.Debug)
	var data []searchapp.Item
	itemSet := make(map[string]int)
	m.s.lock.Lock()
	defer m.s.lock.Unlock()
	for _, rec := range m.s.cliRecords {
		itm, err := searchapp.NewItemFromRecordForQuery(rec, query, m.config.Debug)
		if err != nil {
			// records didn't match the query
			// sugar.Println(" * continue (no match)", rec.Pwd)
			continue
		}
		if idx, ok := itemSet[itm.Key]; ok {
			// duplicate found
			if data[idx].Score >= itm.Score {
				// skip duplicate item
				continue
			}
			// update duplicate item
			data[idx] = itm
			continue
		}
		// add new item
		itemSet[itm.Key] = len(data)
		data = append(data, itm)
	}
	sugar.Debugw("Got new items from records for query, sorting items ...",
		"itemCount", len(data),
	)
	sort.SliceStable(data, func(p, q int) bool {
		return data[p].Score > data[q].Score
	})
	m.s.data = nil
	for _, itm := range data {
		if len(m.s.data) > 420 {
			break
		}
		m.s.data = append(m.s.data, itm)
	}
	m.s.highlightedItem = 0
	sugar.Debugw("Done with data update",
		"recordCount", len(m.s.cliRecords),
		"itemCount", len(m.s.data),
	)
}

func (m manager) UpdateRawData(input string) {
	sugar := m.out.Logger.Sugar()
	sugar.Debugw("Starting RAW data update ...",
		"recordCount", len(m.s.cliRecords),
		"itemCount", len(m.s.data),
	)
	query := searchapp.GetRawTermsFromString(input, m.config.Debug)
	var data []searchapp.RawItem
	itemSet := make(map[string]bool)
	m.s.lock.Lock()
	defer m.s.lock.Unlock()
	for _, rec := range m.s.cliRecords {
		itm, err := searchapp.NewRawItemFromRecordForQuery(rec, query, m.config.Debug)
		if err != nil {
			// records didn't match the query
			// sugar.Println(" * continue (no match)", rec.Pwd)
			continue
		}
		if itemSet[itm.Key] {
			// sugar.Println(" * continue (already present)", itm.key(), itm.pwd)
			continue
		}
		itemSet[itm.Key] = true
		data = append(data, itm)
		// sugar.Println("DATA =", itm.display)
	}
	sugar.Debugw("Got new RAW items from records for query, sorting items ...",
		"itemCount", len(data),
	)
	sort.SliceStable(data, func(p, q int) bool {
		return data[p].Score > data[q].Score
	})
	m.s.rawData = nil
	for _, itm := range data {
		if len(m.s.rawData) > 420 {
			break
		}
		m.s.rawData = append(m.s.rawData, itm)
	}
	m.s.highlightedItem = 0
	sugar.Debugw("Done with RAW data update",
		"recordCount", len(m.s.cliRecords),
		"itemCount", len(m.s.data),
	)
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
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		m.out.Fatal("Failed to set view 'input'", err)
	}

	v.Editable = true
	v.Editor = m
	if m.s.rawMode {
		v.Title = " RESH SEARCH - NON-CONTEXTUAL \"RAW\" MODE - (CTRL+R to switch BACK) "
	} else {
		v.Title = " RESH SEARCH - CONTEXTUAL MODE - (CTRL+R to switch to RAW MODE) "
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
	if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
		m.out.Fatal("Failed to set view 'body'", err)
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

const smallTerminalTresholdWidth = 110

func (m manager) normalMode(g *gocui.Gui, v *gocui.View) error {
	sugar := m.out.Logger.Sugar()
	maxX, maxY := g.Size()

	compactRenderingMode := false
	if maxX < smallTerminalTresholdWidth {
		compactRenderingMode = true
	}

	data := []searchapp.ItemColumns{}

	header := searchapp.GetHeader(compactRenderingMode)
	longestDateLen := len(header.Date)
	longestLocationLen := len(header.Host) + 1 + len(header.PwdTilde)
	longestFlagsLen := 2
	maxPossibleMainViewHeight := maxY - 3 - 1 - 1 - 1 // - top box - header - status - help
	for i, itm := range m.s.data {
		if i == maxY {
			break
		}
		ic := itm.DrawItemColumns(compactRenderingMode, m.config.Debug)
		data = append(data, ic)
		if i > maxPossibleMainViewHeight {
			// do not stretch columns because of results that will end up outside of the page
			continue
		}
		if len(ic.Date) > longestDateLen {
			longestDateLen = len(ic.Date)
		}
		if len(ic.Host)+len(ic.PwdTilde) > longestLocationLen {
			longestLocationLen = len(ic.Host) + len(ic.PwdTilde)
		}
		if len(ic.Flags) > longestFlagsLen {
			longestFlagsLen = len(ic.Flags)
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
	statusLine := searchapp.GetEmptyStatusLine(printedLineLength, realLineLength)
	if m.s.highlightedItem != -1 && m.s.highlightedItem < len(m.s.data) {
		statusLine = m.s.data[m.s.highlightedItem].DrawStatusLine(compactRenderingMode, printedLineLength, realLineLength)
	}
	var statusLineHeight int = len(statusLine)

	helpLineHeight := 1
	const helpLine = "HELP: type to search, UP/DOWN or CTRL+P/N to select, RIGHT to edit, ENTER to execute, CTRL+G to abort, CTRL+C/D to quit; " +
		"FLAGS: G = this git repo, E# = exit status #"
		// "TIP: when resh-cli is launched command line is used as initial search query"

	mainViewHeight := maxY - topBoxHeight - statusLineHeight - helpLineHeight
	m.s.displayedItemsCount = mainViewHeight

	// header
	// header := getHeader()
	// error is expected for header
	dispStr, _, _ := header.ProduceLine(longestDateLen, longestLocationLen, longestFlagsLen, true, true, m.config.Debug)
	dispStr = searchapp.DoHighlightHeader(dispStr, maxX*2)
	v.WriteString(dispStr + "\n")

	var index int
	for index < len(data) {
		itm := data[index]
		if index >= mainViewHeight {
			sugar.Debugw("Reached bottom of the page while producing lines",
				"mainViewHeight", mainViewHeight,
				"predictedMaxViewHeight", maxPossibleMainViewHeight,
			)
			// page is full
			break
		}

		displayStr, _, err := itm.ProduceLine(longestDateLen, longestLocationLen, longestFlagsLen, false, true, m.config.Debug)
		if err != nil {
			sugar.Error("Error while drawing item", zap.Error(err))
		}
		if m.s.highlightedItem == index {
			// maxX * 2 because there are escape sequences that make it hard to tell the real string length
			displayStr = searchapp.DoHighlightString(displayStr, maxX*3)
		}
		if strings.Contains(displayStr, "\n") {
			displayStr = strings.ReplaceAll(displayStr, "\n", "#")
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
	sugar.Debugw("Done drawing page",
		"itemCount", len(m.s.data),
		"highlightedItemIndex", m.s.highlightedItem,
	)
	return nil
}

func (m manager) rawMode(g *gocui.Gui, v *gocui.View) error {
	sugar := m.out.Logger.Sugar()
	maxX, maxY := g.Size()
	topBoxSize := 3
	m.s.displayedItemsCount = maxY - topBoxSize

	for i, itm := range m.s.rawData {
		if i == maxY {
			break
		}
		displayStr := itm.CmdLineWithColor
		if m.s.highlightedItem == i {
			// use actual min requried length instead of 420 constant
			displayStr = searchapp.DoHighlightString(displayStr, maxX*2)
		}
		if strings.Contains(displayStr, "\n") {
			displayStr = strings.ReplaceAll(displayStr, "\n", "#")
		}
		v.WriteString(displayStr + "\n")
	}
	sugar.Debugw("Done drawing page in RAW mode",
		"itemCount", len(m.s.data),
		"highlightedItemIndex", m.s.highlightedItem,
	)
	return nil
}

// SendCliMsg to daemon
func SendCliMsg(out *output.Output, m msg.CliMsg, port string) msg.CliResponse {
	sugar := out.Logger.Sugar()
	recJSON, err := json.Marshal(m)
	if err != nil {
		out.Fatal("Failed to marshal message", err)
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:"+port+"/dump",
		bytes.NewBuffer(recJSON))
	if err != nil {
		out.Fatal("Failed to build request", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		out.FatalDaemonNotRunning(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		out.Fatal("Failed read response", err)
	}
	// sugar.Println(string(body))
	response := msg.CliResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		out.Fatal("Failed decode response", err)
	}
	sugar.Debugw("Recieved records from daemon",
		"recordCount", len(response.CliRecords),
	)
	return response
}
