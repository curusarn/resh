package histfile

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/curusarn/resh/pkg/records"
)

type histfile struct {
	mutex      sync.Mutex
	sessions   map[string]records.Record
	outputPath string
}

// Go creates histfile and runs two gorutines on it
func Go(input chan records.Record, outputPath string, sessionsToDrop chan string) {
	hf := histfile{sessions: map[string]records.Record{}, outputPath: outputPath}
	go hf.writer(input)
	go hf.sessionGC(sessionsToDrop)
}

// sessionGC reads sessionIDs from channel and deletes them from histfile struct
func (h *histfile) sessionGC(sessionsToDrop chan string) {
	for {
		func() {
			session := <-sessionsToDrop
			log.Println("histfile: got session to drop", session)
			h.mutex.Lock()
			defer h.mutex.Unlock()
			if part1, found := h.sessions[session]; found == true {
				log.Println("histfile: Dropping session:", session)
				delete(h.sessions, session)
				go writeRecord(part1, h.outputPath)
			} else {
				log.Println("histfile: No hanging parts for session:", session)
			}
		}()
	}
}

// writer reads records from channel, merges them and writes them to file
func (h *histfile) writer(input chan records.Record) {
	for {
		func() {
			record := <-input
			h.mutex.Lock()
			defer h.mutex.Unlock()

			if record.PartOne {
				if _, found := h.sessions[record.SessionID]; found {
					log.Println("histfile ERROR: Got another first part of the records before merging the previous one - overwriting!")
				}
				h.sessions[record.SessionID] = record
			} else {
				part1, found := h.sessions[record.SessionID]
				if found == false {
					log.Println("histfile ERROR: Got second part of records and nothing to merge it with - ignoring!")
				} else {
					delete(h.sessions, record.SessionID)
					go mergeAndWriteRecord(part1, record, h.outputPath)
				}
			}
		}()
	}
}

func mergeAndWriteRecord(part1, part2 records.Record, outputPath string) {
	err := part1.Merge(part2)
	if err != nil {
		log.Println("Error while merging", err)
		return
	}
	writeRecord(part1, outputPath)
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
