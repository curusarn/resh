package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/curusarn/resh/pkg/records"
)

// HistfileWriter - reads records from channel, merges them and wrotes them to file
func HistfileWriter(input chan records.Record, outputPath string) {
	sessions := map[string]records.Record{}

	for {
		record := <-input
		if record.PartOne {
			if _, found := sessions[record.SessionID]; found {
				log.Println("ERROR: Got another first part of the records before merging the previous one - overwriting!")
			}
			sessions[record.SessionID] = record
		} else {
			part1, found := sessions[record.SessionID]
			if found == false {
				log.Println("ERROR: Got second part of records and nothing to merge it with - ignoring!")
				continue
			}
			delete(sessions, record.SessionID)
			go mergeAndWriteRecord(part1, record, outputPath)
		}
	}
}

func mergeAndWriteRecord(part1, part2 records.Record, outputPath string) {
	err := part1.Merge(part2)
	if err != nil {
		log.Println("Error while merging", err)
		return
	}
	recJSON, err := json.Marshal(part1)
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
		log.Printf("Error while writing: %v, %s\n", part1, err)
		return
	}
}
