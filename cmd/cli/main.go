package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/awesome-gocui/gocui"
	"github.com/curusarn/resh/pkg/cfg"
	"github.com/curusarn/resh/pkg/msg"

	"os/user"
	"path/filepath"
	"strconv"
)

// version from git set during build
var version string

// commit from git set during build
var commit string

func main() {
	usr, _ := user.Current()
	dir := usr.HomeDir
	configPath := filepath.Join(dir, "/.config/resh.toml")

	var config cfg.Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Fatal("Error reading config:", err)
	}

	sessionID := flag.String("sessionID", "", "resh generated session id")
	flag.Parse()

	if *sessionID == "" {
		fmt.Println("Error: you need to specify sessionId")
	}

	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	// g.SelBgColor = gocui.ColorGreen
	g.Highlight = false

	mess := msg.InspectMsg{SessionID: *sessionID, Count: 40}
	resp := SendInspectMsg(mess, strconv.Itoa(config.Port))

	st := state{
		// lock sync.Mutex
		dataOriginal: resp.CmdLines,
		data:         resp.CmdLines,
	}
	layout := manager{
		sessionID: *sessionID,
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

	err = g.MainLoop()
	if err != nil && gocui.IsQuit(err) == false {
		log.Panicln(err)
	}
	layout.Output()
}

// returns the number of hits for query
func queryHits(cmdline string, queryTerms []string) int {
	hits := 0
	for _, term := range queryTerms {
		if strings.Contains(cmdline, term) {
			hits++
		}
	}
	return hits
}

type state struct {
	dataOriginal    []string
	data            []string
	highlightedItem int

	outputBuffer string
}

type manager struct {
	sessionID string
	config    cfg.Config

	s *state
}

func (m manager) Output() {
	if len(m.s.outputBuffer) > 0 {
		fmt.Print(m.s.outputBuffer)
	}
}

func (m manager) SelectExecute(g *gocui.Gui, v *gocui.View) error {
	if m.s.highlightedItem < len(m.s.data) {
		m.s.outputBuffer = m.s.data[m.s.highlightedItem]
		return gocui.ErrQuit
	}
	return nil
}

func (m manager) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	gocui.DefaultEditor.Edit(v, key, ch, mod)
	query := v.Buffer()
	terms := strings.Split(query, " ")
	var dataHits []int
	m.s.data = nil
	for _, entry := range m.s.dataOriginal {
		hits := queryHits(entry, terms)
		if hits > 0 {
			m.s.data = append(m.s.data, entry)
			dataHits = append(dataHits, hits)
		}
	}
	sort.SliceStable(m.s.data, func(p, q int) bool {
		return dataHits[p] > dataHits[q]
	})
	m.s.highlightedItem = 0
}

func (m manager) Next(g *gocui.Gui, v *gocui.View) error {
	_, y := g.Size()
	if m.s.highlightedItem < y {
		m.s.highlightedItem++
	}
	return nil
}

func (m manager) Prev(g *gocui.Gui, v *gocui.View) error {
	if m.s.highlightedItem > 0 {
		m.s.highlightedItem--
	}
	return nil
}

// you can have Layout with pointer reciever if you pass the layout function to the setmanger
// I dont think we need that tho
func (m manager) Layout(g *gocui.Gui) error {
	var b byte
	maxX, maxY := g.Size()

	v, err := g.SetView("input", 0, 0, maxX-1, 2, b)
	if err != nil && gocui.IsUnknownView(err) == false {
		log.Panicln(err.Error())
	}

	v.Editable = true
	// v.Editor = gocui.EditorFunc(m.editor.Edit)
	v.Editor = m
	v.Title = "resh cli"

	g.SetCurrentView("input")

	v, err = g.SetView("body", 0, 2, maxX-1, maxY, b)
	if err != nil && gocui.IsUnknownView(err) == false {
		log.Panicln(err.Error())
	}
	v.Frame = false
	v.Autoscroll = true
	v.Clear()
	for _, cmdLine := range m.s.data {
		entry := strings.Trim(cmdLine, "\n") + "\n"
		v.WriteString(entry)
	}
	if m.s.highlightedItem < len(m.s.data) {
		v.SetHighlight(m.s.highlightedItem, true)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// SendInspectMsg to daemon
func SendInspectMsg(m msg.InspectMsg, port string) msg.MultiResponse {
	recJSON, err := json.Marshal(m)
	if err != nil {
		log.Fatal("send err 1", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:"+port+"/inspect",
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
	response := msg.MultiResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal("unmarshal resp error: ", err)
	}
	return response
}
