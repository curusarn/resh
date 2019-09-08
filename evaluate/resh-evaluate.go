package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
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

	showVersion := flag.Bool("version", false, "Show version and exit")
	showRevision := flag.Bool("revision", false, "Show git revision and exit")
	inputPath := flag.String("input", "",
		"Input file (default: "+historyPath+"OR"+sanitizedHistoryPath+
			" depending on --sanitized-input option)")
	outputPath := flag.String("output", "", "Output file (default: use stdout)")
	sanitizedInput := flag.Bool("sanitized-input", false,
		"Handle input as sanitized (also changes default value for input argument)")

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

	var writer *bufio.Writer
	if *outputPath != "" {
		outputFile, err := os.Create(*outputPath)
		if err != nil {
			log.Fatal("Create() output file error:", err)
		}
		defer outputFile.Close()
		writer = bufio.NewWriter(outputFile)
	} else {
		writer = bufio.NewWriter(os.Stdout)
	}
	defer writer.Flush()

	evaluator := evaluator{sanitizedInput: *sanitizedInput, writer: writer, maxCandidates: 42}
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
}

type strategy interface {
	GetTitleAndDescription() (string, string)
	GetCandidates() []string
	AddHistoryRecord(record *common.Record) error
	ResetHistory() error
}

type evaluator struct {
	sanitizedInput bool
	writer         *bufio.Writer
	maxCandidates  int
	historyRecords []common.Record
}

func (e *evaluator) init(inputPath string) error {
	e.historyRecords = e.loadHistoryRecords(inputPath)
	return nil
}

func (e *evaluator) evaluate(strat strategy) error {
	stats := statistics{writer: e.writer, size: e.maxCandidates + 1}
	stats.init()

	for _, record := range e.historyRecords {
		candidates := strat.GetCandidates()

		match := false
		for i, candidate := range candidates {
			// make an option (--calculate-total) to turn this on/off ?
			// if i >= e.maxCandidates {
			// 	break
			// }
			if candidate == record.CmdLine {
				stats.addMatch(i+1, record.CmdLength)
				match = true
				break
			}
		}
		if match == false {
			stats.addMiss()
		}
		err := strat.AddHistoryRecord(&record)
		if err != nil {
			log.Println("Error while evauating", err)
			return err
		}
	}
	title, description := strat.GetTitleAndDescription()
	n, err := e.writer.WriteString(title + " - " + description + "\n")
	if err != nil {
		log.Fatal(err)
	}
	if n == 0 {
		log.Fatal("Nothing was written", n)
	}
	// print results
	stats.printCumulative()
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
