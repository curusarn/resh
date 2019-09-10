package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/curusarn/resh/common"
)

// Version from git set during build
var Version string

// Revision from git set during build
var Revision string

func main() {
	usr, _ := user.Current()
	dir := usr.HomeDir
	historyPath := filepath.Join(dir, ".resh_history.json")
	sanitizedHistoryPath := filepath.Join(dir, "resh_history_sanitized.json")
	// tmpPath := "/tmp/resh-evaluate-tmp.json"

	showVersion := flag.Bool("version", false, "Show version and exit")
	showRevision := flag.Bool("revision", false, "Show git revision and exit")
	inputPath := flag.String("input", "",
		"Input file (default: "+historyPath+"OR"+sanitizedHistoryPath+
			" depending on --sanitized-input option)")
	// outputDir := flag.String("output", "/tmp/resh-evaluate", "Output directory")
	sanitizedInput := flag.Bool("sanitized-input", false,
		"Handle input as sanitized (also changes default value for input argument)")
	plottingScript := flag.String("plotting-script", "resh-evaluate-plot.py", "Script to use for plotting")

	flag.Parse()

	// set default input
	if *inputPath == "" {
		if *sanitizedInput {
			*inputPath = sanitizedHistoryPath
		} else {
			*inputPath = historyPath
		}
	}

	if *showVersion == true {
		fmt.Println(Version)
		os.Exit(0)
	}
	if *showRevision == true {
		fmt.Println(Revision)
		os.Exit(0)
	}

	evaluator := evaluator{sanitizedInput: *sanitizedInput, maxCandidates: 50}
	err := evaluator.init(*inputPath)
	if err != nil {
		log.Fatal("Evaluator init() error:", err)
	}

	var strategies []strategy

	// dummy := strategyDummy{}
	// strategies = append(strategies, &dummy)

	recent := strategyRecent{}
	strategies = append(strategies, &recent)

	for _, strat := range strategies {
		err = evaluator.evaluate(strat)
		if err != nil {
			log.Println("Evaluator evaluate() error:", err)
		}
	}
	// evaluator.dumpJSON(tmpPath)

	evaluator.calculateStatsAndPlot(*plottingScript)
}

type strategy interface {
	GetTitleAndDescription() (string, string)
	GetCandidates() []string
	AddHistoryRecord(record *common.Record) error
	ResetHistory() error
}

type matchJSON struct {
	Match         bool
	Distance      int
	CharsRecalled int
}

type strategyJSON struct {
	Title       string
	Description string
	Matches     []matchJSON
}

type evaluateJSON struct {
	Strategies []strategyJSON
	Records    []common.Record
}

type evaluator struct {
	sanitizedInput bool
	maxCandidates  int
	historyRecords []common.Record
	data           evaluateJSON
}

func (e *evaluator) init(inputPath string) error {
	e.historyRecords = e.loadHistoryRecords(inputPath)
	e.processRecords()
	return nil
}

func (e *evaluator) calculateStatsAndPlot(scriptName string) {
	evalJSON, err := json.Marshal(e.data)
	if err != nil {
		log.Fatal("json marshal error", err)
	}
	buffer := bytes.Buffer{}
	buffer.Write(evalJSON)
	// run python script to stat and plot/
	cmd := exec.Command(scriptName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = &buffer
	err = cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}
}

// enrich records and add them to serializable structure
func (e *evaluator) processRecords() {
	for _, record := range e.historyRecords {

		// assert
		if record.Sanitized != e.sanitizedInput {
			if e.sanitizedInput {
				log.Fatal("ASSERT failed: '--sanitized-input' is present but data is not sanitized")
			}
			log.Fatal("ASSERT failed: data is sanitized but '--sanitized-input' is not present")
		}

		record.Enrich()
		e.data.Records = append(e.data.Records, record)
	}
}

func (e *evaluator) evaluate(strategy strategy) error {
	title, description := strategy.GetTitleAndDescription()
	strategyData := strategyJSON{Title: title, Description: description}
	for _, record := range e.historyRecords {
		candidates := strategy.GetCandidates()

		matchFound := false
		for i, candidate := range candidates {
			// make an option (--calculate-total) to turn this on/off ?
			// if i >= e.maxCandidates {
			// 	break
			// }
			if candidate == record.CmdLine {
				match := matchJSON{Match: true, Distance: i + 1, CharsRecalled: record.CmdLength}
				strategyData.Matches = append(strategyData.Matches, match)
				matchFound = true
				break
			}
		}
		if matchFound == false {
			strategyData.Matches = append(strategyData.Matches, matchJSON{})
		}
		err := strategy.AddHistoryRecord(&record)
		if err != nil {
			log.Println("Error while evauating", err)
			return err
		}
	}
	e.data.Strategies = append(e.data.Strategies, strategyData)
	return nil
}

func (e *evaluator) loadHistoryRecords(fname string) []common.Record {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal("Open() resh history file error:", err)
	}
	defer file.Close()

	var records []common.Record
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := common.Record{}
		fallbackRecord := common.FallbackRecord{}
		line := scanner.Text()
		err = json.Unmarshal([]byte(line), &record)
		if err != nil {
			err = json.Unmarshal([]byte(line), &fallbackRecord)
			if err != nil {
				log.Println("Line:", line)
				log.Fatal("Decoding error:", err)
			}
			record = common.ConvertRecord(&fallbackRecord)
		}
		if e.sanitizedInput == false {
			if record.CmdLength != 0 {
				log.Fatal("Assert failed - 'cmdLength' is set in raw data. Maybe you want to use '--sanitized-input' option?")
			}
			record.CmdLength = len(record.CmdLine)
		}
		if record.CmdLength == 0 {
			log.Fatal("Assert failed - 'cmdLength' is unset in the data. This should not happen.")
		}
		records = append(records, record)
	}
	return records
}
