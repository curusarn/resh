package records

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"testing"
)

func GetTestRecords() []Record {
	file, err := os.Open("testdata/resh_history.json")
	if err != nil {
		log.Fatalf("Failed to open resh history file: %v", err)
	}
	defer file.Close()

	var recs []Record
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := Record{}
		line := scanner.Text()
		err = json.Unmarshal([]byte(line), &record)
		if err != nil {
			log.Fatalf("Error decoding record: '%s'; err: %v", line, err)
		}
		recs = append(recs, record)
	}
	return recs
}

func GetTestEnrichedRecords() []EnrichedRecord {
	var recs []EnrichedRecord
	for _, rec := range GetTestRecords() {
		recs = append(recs, Enriched(rec))
	}
	return recs
}

func TestToString(t *testing.T) {
	for _, rec := range GetTestEnrichedRecords() {
		_, err := rec.ToString()
		if err != nil {
			t.Error("ToString() failed")
		}
	}
}

func TestEnriched(t *testing.T) {
	record := Record{BaseRecord: BaseRecord{CmdLine: "cmd arg1 arg2"}}
	enriched := Enriched(record)
	if enriched.FirstWord != "cmd" || enriched.Command != "cmd" {
		t.Error("Enriched() returned reocord w/ wrong Command OR FirstWord")
	}
}

func TestValidate(t *testing.T) {
	record := EnrichedRecord{}
	if record.Validate() == nil {
		t.Error("Validate() didn't return an error for invalid record")
	}
	record.CmdLine = "cmd arg"
	record.FirstWord = "cmd"
	record.Command = "cmd"
	time := 1234.5678
	record.RealtimeBefore = time
	record.RealtimeAfter = time
	record.RealtimeBeforeLocal = time
	record.RealtimeAfterLocal = time
	pwd := "/pwd"
	record.Pwd = pwd
	record.PwdAfter = pwd
	record.RealPwd = pwd
	record.RealPwdAfter = pwd
	if record.Validate() != nil {
		t.Error("Validate() returned an error for a valid record")
	}
}

func TestGetCommandAndFirstWord(t *testing.T) {
	cmd, stWord, err := GetCommandAndFirstWord("cmd arg1 arg2")
	if err != nil || cmd != "cmd" || stWord != "cmd" {
		t.Error("GetCommandAndFirstWord() returned wrong Command OR FirstWord")
	}
}
