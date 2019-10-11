package histanal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/curusarn/resh/pkg/records"
)

type deviceRecords struct {
	Name    string
	Records []records.EnrichedRecord
}

type userRecords struct {
	Name    string
	Devices []deviceRecords
}

// HistLoad loads history
type HistLoad struct {
	UsersRecords   []userRecords
	skipFailedCmds bool
	sanitizedInput bool
	debugRecords   float64
}

func (e *HistLoad) preprocessDeviceRecords(device deviceRecords) deviceRecords {
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
func (e *HistLoad) preprocessRecords() {
	for i := range e.UsersRecords {
		for j := range e.UsersRecords[i].Devices {
			e.UsersRecords[i].Devices[j] = e.preprocessDeviceRecords(e.UsersRecords[i].Devices[j])
		}
	}
}

func (e *HistLoad) loadHistoryRecordsBatchMode(fname string, dataRootPath string) []userRecords {
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

func (e *HistLoad) loadHistoryRecords(fname string) []records.EnrichedRecord {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal("Open() resh history file error:", err)
	}
	defer file.Close()

	var recs []records.EnrichedRecord
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := records.Record{}
		fallbackRecord := records.FallbackRecord{}
		line := scanner.Text()
		err = json.Unmarshal([]byte(line), &record)
		if err != nil {
			err = json.Unmarshal([]byte(line), &fallbackRecord)
			if err != nil {
				log.Println("Line:", line)
				log.Fatal("Decoding error:", err)
			}
			record = records.Convert(&fallbackRecord)
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
		recs = append(recs, records.Enriched(record))
	}
	return recs
}
