package histfile

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/curusarn/resh/pkg/records"
)

// Histfile writes records to histfile
type Histfile struct {
	sessionsMutex sync.Mutex
	sessions      map[string]records.Record
	historyPath   string

	recentMutex       sync.Mutex
	recentRecords     []records.Record
	recentCmdLines    []string // deduplicated
	cmdLinesLastIndex map[string]int
}

// New creates new histfile and runs two gorutines on it
func New(input chan records.Record, historyPath string, initHistSize int, sessionsToDrop chan string) *Histfile {
	hf := Histfile{
		sessions:          map[string]records.Record{},
		historyPath:       historyPath,
		cmdLinesLastIndex: map[string]int{},
	}
	go hf.loadHistory(initHistSize)
	go hf.writer(input)
	go hf.sessionGC(sessionsToDrop)
	return &hf
}

func (h *Histfile) loadHistory(initHistSize int) {
	h.recentMutex.Lock()
	defer h.recentMutex.Unlock()
	h.recentCmdLines = records.LoadCmdLinesFromFile(h.historyPath, initHistSize)
}

// sessionGC reads sessionIDs from channel and deletes them from histfile struct
func (h *Histfile) sessionGC(sessionsToDrop chan string) {
	for {
		func() {
			session := <-sessionsToDrop
			log.Println("histfile: got session to drop", session)
			h.sessionsMutex.Lock()
			defer h.sessionsMutex.Unlock()
			if part1, found := h.sessions[session]; found == true {
				log.Println("histfile: Dropping session:", session)
				delete(h.sessions, session)
				go writeRecord(part1, h.historyPath)
			} else {
				log.Println("histfile: No hanging parts for session:", session)
			}
		}()
	}
}

// writer reads records from channel, merges them and writes them to file
func (h *Histfile) writer(input chan records.Record) {
	for {
		func() {
			record := <-input
			h.sessionsMutex.Lock()
			defer h.sessionsMutex.Unlock()

			if record.PartOne {
				if _, found := h.sessions[record.SessionID]; found {
					log.Println("histfile WARN: Got another first part of the records before merging the previous one - overwriting! " +
						"(this happens in bash because bash-preexec runs when it's not supposed to)")
				}
				h.sessions[record.SessionID] = record
			} else {
				if part1, found := h.sessions[record.SessionID]; found == false {
					log.Println("histfile ERROR: Got second part of records and nothing to merge it with - ignoring!")
				} else {
					delete(h.sessions, record.SessionID)
					go h.mergeAndWriteRecord(part1, record)
				}
			}
		}()
	}
}

func (h *Histfile) mergeAndWriteRecord(part1, part2 records.Record) {
	err := part1.Merge(part2)
	if err != nil {
		log.Println("Error while merging", err)
		return
	}

	func() {
		h.recentMutex.Lock()
		defer h.recentMutex.Unlock()
		h.recentRecords = append(h.recentRecords, part1)
		cmdLine := part1.CmdLine
		idx, found := h.cmdLinesLastIndex[cmdLine]
		if found {
			h.recentCmdLines = append(h.recentCmdLines[:idx], h.recentCmdLines[idx+1:]...)
		}
		h.cmdLinesLastIndex[cmdLine] = len(h.recentCmdLines)
		h.recentCmdLines = append(h.recentCmdLines, cmdLine)
	}()

	writeRecord(part1, h.historyPath)
}

func writeRecord(rec records.Record, outputPath string) {
	recJSON, err := json.Marshal(rec)
	if err != nil {
		log.Println("Marshalling error", err)
		return
	}
	f, err := os.OpenFile(outputPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Could not open file", err)
		return
	}
	defer f.Close()
	_, err = f.Write(append(recJSON, []byte("\n")...))
	if err != nil {
		log.Printf("Error while writing: %v, %s\n", rec, err)
		return
	}
}

// GetRecentCmdLines returns recent cmdLines
func (h *Histfile) GetRecentCmdLines(limit int) []string {
	return h.recentCmdLines
}
