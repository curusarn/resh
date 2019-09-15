package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sort"

	"github.com/curusarn/resh/common"
	"github.com/jpillora/longestcommon"
	"github.com/schollz/progressbar"
)

// Version from git set during build
var Version string

// Revision from git set during build
var Revision string

func main() {
	const maxCandidates = 50

	usr, _ := user.Current()
	dir := usr.HomeDir
	historyPath := filepath.Join(dir, ".resh_history.json")
	historyPathBatchMode := filepath.Join(dir, "resh_history.json")
	sanitizedHistoryPath := filepath.Join(dir, "resh_history_sanitized.json")
	// tmpPath := "/tmp/resh-evaluate-tmp.json"

	showVersion := flag.Bool("version", false, "Show version and exit")
	showRevision := flag.Bool("revision", false, "Show git revision and exit")
	input := flag.String("input", "",
		"Input file (default: "+historyPath+"OR"+sanitizedHistoryPath+
			" depending on --sanitized-input option)")
	// outputDir := flag.String("output", "/tmp/resh-evaluate", "Output directory")
	sanitizedInput := flag.Bool("sanitized-input", false,
		"Handle input as sanitized (also changes default value for input argument)")
	plottingScript := flag.String("plotting-script", "resh-evaluate-plot.py", "Script to use for plotting")
	inputDataRoot := flag.String("input-data-root", "",
		"Input data root, enables batch mode, looks for files matching --input option")
	slow := flag.Bool("slow", false,
		"Enables stuff that takes a long time (e.g. markov chain strategies).")

	flag.Parse()

	// handle show{Version,Revision} options
	if *showVersion == true {
		fmt.Println(Version)
		os.Exit(0)
	}
	if *showRevision == true {
		fmt.Println(Revision)
		os.Exit(0)
	}

	// handle batch mode
	batchMode := false
	if *inputDataRoot != "" {
		batchMode = true
	}
	// set default input
	if *input == "" {
		if *sanitizedInput {
			*input = sanitizedHistoryPath
		} else if batchMode {
			*input = historyPathBatchMode
		} else {
			*input = historyPath
		}
	}

	evaluator := evaluator{sanitizedInput: *sanitizedInput, maxCandidates: maxCandidates, BatchMode: batchMode}
	if batchMode {
		err := evaluator.initBatchMode(*input, *inputDataRoot)
		if err != nil {
			log.Fatal("Evaluator initBatchMode() error:", err)
		}
	} else {
		err := evaluator.init(*input)
		if err != nil {
			log.Fatal("Evaluator init() error:", err)
		}
	}

	var strategies []strategy

	// dummy := strategyDummy{}
	// strategies = append(strategies, &dummy)

	recent := strategyRecent{}
	frequent := strategyFrequent{}
	frequent.init()
	directory := strategyDirectorySensitive{}
	directory.init()
	random := strategyRandom{candidatesSize: maxCandidates}
	random.init()

	markovCmd := strategyMarkovChainCmd{order: 1}
	markovCmd.init()

	markovCmd2 := strategyMarkovChainCmd{order: 2}
	markovCmd2.init()

	markov := strategyMarkovChain{order: 1}
	markov.init()

	markov2 := strategyMarkovChain{order: 2}
	markov2.init()

	strategies = append(strategies, &recent, &frequent, &directory, &random)

	if *slow {
		strategies = append(strategies, &markovCmd2, &markovCmd, &markov2, &markov)
	}

	for _, strat := range strategies {
		err := evaluator.evaluate(strat)
		if err != nil {
			log.Println("Evaluator evaluate() error:", err)
		}
	}

	evaluator.calculateStatsAndPlot(*plottingScript)
}

type strategy interface {
	GetTitleAndDescription() (string, string)
	GetCandidates() []string
	AddHistoryRecord(record *common.EnrichedRecord) error
	ResetHistory() error
}

type matchJSON struct {
	Match         bool
	Distance      int
	CharsRecalled int
}

type multiMatchItemJSON struct {
	Distance      int
	CharsRecalled int
}

type multiMatchJSON struct {
	Match   bool
	Entries []multiMatchItemJSON
}

type strategyJSON struct {
	Title         string
	Description   string
	Matches       []matchJSON
	PrefixMatches []multiMatchJSON
}

type deviceRecords struct {
	Name    string
	Records []common.EnrichedRecord
}

type userRecords struct {
	Name    string
	Devices []deviceRecords
}

type evaluator struct {
	sanitizedInput bool
	BatchMode      bool
	maxCandidates  int
	UsersRecords   []userRecords
	Strategies     []strategyJSON
}

func (e *evaluator) initBatchMode(input string, inputDataRoot string) error {
	e.UsersRecords = e.loadHistoryRecordsBatchMode(input, inputDataRoot)
	e.processRecords()
	return nil
}

func (e *evaluator) init(inputPath string) error {
	records := e.loadHistoryRecords(inputPath)
	device := deviceRecords{Records: records}
	user := userRecords{}
	user.Devices = append(user.Devices, device)
	e.UsersRecords = append(e.UsersRecords, user)
	e.processRecords()
	return nil
}

func (e *evaluator) calculateStatsAndPlot(scriptName string) {
	evalJSON, err := json.Marshal(e)
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
	for i := range e.UsersRecords {
		for j, device := range e.UsersRecords[i].Devices {
			sessionIDs := map[string]uint64{}
			var nextID uint64
			nextID = 1 // start with 1 because 0 won't get saved to json
			for k, record := range e.UsersRecords[i].Devices[j].Records {
				id, found := sessionIDs[record.SessionID]
				if found == false {
					id = nextID
					sessionIDs[record.SessionID] = id
					nextID++
				}
				e.UsersRecords[i].Devices[j].Records[k].SeqSessionID = id
				// assert
				if record.Sanitized != e.sanitizedInput {
					if e.sanitizedInput {
						log.Fatal("ASSERT failed: '--sanitized-input' is present but data is not sanitized")
					}
					log.Fatal("ASSERT failed: data is sanitized but '--sanitized-input' is not present")
				}
			}
			sort.SliceStable(e.UsersRecords[i].Devices[j].Records, func(x, y int) bool {
				if device.Records[x].SeqSessionID == device.Records[y].SeqSessionID {
					return device.Records[x].RealtimeAfterLocal < device.Records[y].RealtimeAfterLocal
				}
				return device.Records[x].SeqSessionID < device.Records[y].SeqSessionID
			})
		}
	}
}

func (e *evaluator) evaluate(strategy strategy) error {
	title, description := strategy.GetTitleAndDescription()
	log.Println("Evaluating strategy:", title, "-", description)
	strategyData := strategyJSON{Title: title, Description: description}
	for i := range e.UsersRecords {
		for j := range e.UsersRecords[i].Devices {
			bar := progressbar.New(len(e.UsersRecords[i].Devices[j].Records))
			for _, record := range e.UsersRecords[i].Devices[j].Records {
				candidates := strategy.GetCandidates()

				matchFound := false
				longestPrefixMatchLength := 0
				multiMatch := multiMatchJSON{}
				for i, candidate := range candidates {
					// make an option (--calculate-total) to turn this on/off ?
					// if i >= e.maxCandidates {
					// 	break
					// }
					commonPrefixLength := len(longestcommon.Prefix([]string{candidate, record.CmdLine}))
					if commonPrefixLength > longestPrefixMatchLength {
						longestPrefixMatchLength = commonPrefixLength
						prefixMatch := multiMatchItemJSON{Distance: i + 1, CharsRecalled: commonPrefixLength}
						multiMatch.Match = true
						multiMatch.Entries = append(multiMatch.Entries, prefixMatch)
					}
					if candidate == record.CmdLine {
						match := matchJSON{Match: true, Distance: i + 1, CharsRecalled: record.CmdLength}
						matchFound = true
						strategyData.Matches = append(strategyData.Matches, match)
						strategyData.PrefixMatches = append(strategyData.PrefixMatches, multiMatch)
						break
					}
				}
				if matchFound == false {
					strategyData.Matches = append(strategyData.Matches, matchJSON{})
					strategyData.PrefixMatches = append(strategyData.PrefixMatches, multiMatch)
				}
				err := strategy.AddHistoryRecord(&record)
				if err != nil {
					log.Println("Error while evauating", err)
					return err
				}
				bar.Add(1)
			}
			strategy.ResetHistory()
			fmt.Println()
		}
	}
	e.Strategies = append(e.Strategies, strategyData)
	return nil
}

func (e *evaluator) loadHistoryRecordsBatchMode(fname string, dataRootPath string) []userRecords {
	var records []userRecords
	info, err := os.Stat(dataRootPath)
	if err != nil {
		log.Fatal("Error: Directory", dataRootPath, "does not exist - exiting! (", err, ")")
	}
	if info.IsDir() == false {
		log.Fatal("Error:", dataRootPath, "is not a directory - exiting!")
	}
	users, err := ioutil.ReadDir(dataRootPath)
	if err != nil {
		log.Fatal("Could not read directory:", dataRootPath)
	}
	fmt.Println("Listing users in <", dataRootPath, ">...")
	for _, user := range users {
		userRecords := userRecords{Name: user.Name()}
		userFullPath := filepath.Join(dataRootPath, user.Name())
		if user.IsDir() == false {
			log.Println("Warn: Unexpected file (not a directory) <", userFullPath, "> - skipping.")
			continue
		}
		fmt.Println()
		fmt.Printf("*- %s\n", user.Name())
		devices, err := ioutil.ReadDir(userFullPath)
		if err != nil {
			log.Fatal("Could not read directory:", userFullPath)
		}
		for _, device := range devices {
			deviceRecords := deviceRecords{Name: device.Name()}
			deviceFullPath := filepath.Join(userFullPath, device.Name())
			if device.IsDir() == false {
				log.Println("Warn: Unexpected file (not a directory) <", deviceFullPath, "> - skipping.")
				continue
			}
			fmt.Printf("   \\- %s\n", device.Name())
			files, err := ioutil.ReadDir(deviceFullPath)
			if err != nil {
				log.Fatal("Could not read directory:", deviceFullPath)
			}
			for _, file := range files {
				fileFullPath := filepath.Join(deviceFullPath, file.Name())
				if file.Name() == fname {
					fmt.Printf("      \\- %s - loading ...", file.Name())
					// load the data
					deviceRecords.Records = e.loadHistoryRecords(fileFullPath)
					fmt.Println(" OK ✓")
				} else {
					fmt.Printf("      \\- %s - skipped\n", file.Name())
				}
			}
			userRecords.Devices = append(userRecords.Devices, deviceRecords)
		}
		records = append(records, userRecords)
	}
	return records
}

func (e *evaluator) loadHistoryRecords(fname string) []common.EnrichedRecord {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal("Open() resh history file error:", err)
	}
	defer file.Close()

	var records []common.EnrichedRecord
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
		records = append(records, record.Enrich())
	}
	return records
}
