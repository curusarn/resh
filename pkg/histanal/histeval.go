package histanal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"

	"github.com/curusarn/resh/pkg/records"
	"github.com/curusarn/resh/pkg/strat"
	"github.com/jpillora/longestcommon"

	"github.com/schollz/progressbar"
)

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

// HistEval evaluates history
type HistEval struct {
	HistLoad
	BatchMode     bool
	maxCandidates int
	Strategies    []strategyJSON
}

// NewHistEval constructs new HistEval
func NewHistEval(inputPath string,
	maxCandidates int, skipFailedCmds bool,
	debugRecords float64, sanitizedInput bool) HistEval {

	e := HistEval{
		HistLoad: HistLoad{
			skipFailedCmds: skipFailedCmds,
			debugRecords:   debugRecords,
			sanitizedInput: sanitizedInput,
		},
		maxCandidates: maxCandidates,
		BatchMode:     false,
	}
	records := e.loadHistoryRecords(inputPath)
	device := deviceRecords{Records: records}
	user := userRecords{}
	user.Devices = append(user.Devices, device)
	e.UsersRecords = append(e.UsersRecords, user)
	e.preprocessRecords()
	return e
}

// NewHistEvalBatchMode constructs new HistEval in batch mode
func NewHistEvalBatchMode(input string, inputDataRoot string,
	maxCandidates int, skipFailedCmds bool,
	debugRecords float64, sanitizedInput bool) HistEval {

	e := HistEval{
		HistLoad: HistLoad{
			skipFailedCmds: skipFailedCmds,
			debugRecords:   debugRecords,
			sanitizedInput: sanitizedInput,
		},
		maxCandidates: maxCandidates,
		BatchMode:     false,
	}
	e.UsersRecords = e.loadHistoryRecordsBatchMode(input, inputDataRoot)
	e.preprocessRecords()
	return e
}

func (e *HistEval) preprocessDeviceRecords(device deviceRecords) deviceRecords {
	sessionIDs := map[string]uint64{}
	var nextID uint64
	nextID = 1 // start with 1 because 0 won't get saved to json
	for k, record := range device.Records {
		id, found := sessionIDs[record.SessionID]
		if found == false {
			id = nextID
			sessionIDs[record.SessionID] = id
			nextID++
		}
		device.Records[k].SeqSessionID = id
		// assert
		if record.Sanitized != e.sanitizedInput {
			if e.sanitizedInput {
				log.Fatal("ASSERT failed: '--sanitized-input' is present but data is not sanitized")
			}
			log.Fatal("ASSERT failed: data is sanitized but '--sanitized-input' is not present")
		}
		device.Records[k].SeqSessionID = id
		if e.debugRecords > 0 && rand.Float64() < e.debugRecords {
			device.Records[k].DebugThisRecord = true
		}
	}
	// sort.SliceStable(device.Records, func(x, y int) bool {
	// 	if device.Records[x].SeqSessionID == device.Records[y].SeqSessionID {
	// 		return device.Records[x].RealtimeAfterLocal < device.Records[y].RealtimeAfterLocal
	// 	}
	// 	return device.Records[x].SeqSessionID < device.Records[y].SeqSessionID
	// })

	// iterate from back and mark last record of each session
	sessionIDSet := map[string]bool{}
	for i := len(device.Records) - 1; i >= 0; i-- {
		var record *records.EnrichedRecord
		record = &device.Records[i]
		if sessionIDSet[record.SessionID] {
			continue
		}
		sessionIDSet[record.SessionID] = true
		record.LastRecordOfSession = true
	}
	return device
}

// enrich records and add sequential session ID
func (e *HistEval) preprocessRecords() {
	for i := range e.UsersRecords {
		for j := range e.UsersRecords[i].Devices {
			e.UsersRecords[i].Devices[j] = e.preprocessDeviceRecords(e.UsersRecords[i].Devices[j])
		}
	}
}

// Evaluate a given strategy
func (e *HistEval) Evaluate(strategy strat.IStrategy) error {
	title, description := strategy.GetTitleAndDescription()
	log.Println("Evaluating strategy:", title, "-", description)
	strategyData := strategyJSON{Title: title, Description: description}
	for i := range e.UsersRecords {
		for j := range e.UsersRecords[i].Devices {
			bar := progressbar.New(len(e.UsersRecords[i].Devices[j].Records))
			var prevRecord records.EnrichedRecord
			for _, record := range e.UsersRecords[i].Devices[j].Records {
				if e.skipFailedCmds && record.ExitCode != 0 {
					continue
				}
				candidates := strategy.GetCandidates(records.Stripped(record))
				if record.DebugThisRecord {
					log.Println()
					log.Println("===================================================")
					log.Println("STRATEGY:", title, "-", description)
					log.Println("===================================================")
					log.Println("Previous record:")
					if prevRecord.RealtimeBefore == 0 {
						log.Println("== NIL")
					} else {
						rec, _ := prevRecord.ToString()
						log.Println(rec)
					}
					log.Println("---------------------------------------------------")
					log.Println("Recommendations for:")
					rec, _ := record.ToString()
					log.Println(rec)
					log.Println("---------------------------------------------------")
					for i, candidate := range candidates {
						if i > 10 {
							break
						}
						log.Println(string(candidate))
					}
					log.Println("===================================================")
				}

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
				prevRecord = record
			}
			strategy.ResetHistory()
			fmt.Println()
		}
	}
	e.Strategies = append(e.Strategies, strategyData)
	return nil
}

// CalculateStatsAndPlot results
func (e *HistEval) CalculateStatsAndPlot(scriptName string) {
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
