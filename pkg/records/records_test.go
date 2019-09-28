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
		log.Fatal("Open() resh history file error:", err)
	}
	defer file.Close()

	var recs []Record
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := Record{}
		line := scanner.Text()
		err = json.Unmarshal([]byte(line), &record)
		if err != nil {
			log.Println("Line:", line)
			log.Fatal("Decoding error:", err)
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

func TestSetCmdLine(t *testing.T) {
	record := EnrichedRecord{}
	cmdline := "cmd arg1 arg2"
	record.SetCmdLine(cmdline)
	if record.CmdLine != cmdline || record.Command != "cmd" || record.FirstWord != "cmd" {
		t.Error()
	}
}

func TestStripped(t *testing.T) {
	for _, rec := range GetTestEnrichedRecords() {
		stripped := Stripped(rec)

		// there should be no cmdline
		if stripped.CmdLine != "" ||
			stripped.FirstWord != "" ||
			stripped.Command != "" {
			t.Error("Stripped() returned record w/ info about CmdLine, Command OR FirstWord")
		}
		//  *after* fields should be overwritten by *before* fields
		if stripped.PwdAfter != stripped.Pwd ||
			stripped.RealPwdAfter != stripped.RealPwd ||
			stripped.TimezoneAfter != stripped.TimezoneBefore ||
			stripped.RealtimeAfter != stripped.RealtimeBefore ||
			stripped.RealtimeAfterLocal != stripped.RealtimeBeforeLocal {
			t.Error("Stripped() returned record w/ different *after* and *before* values - *after* fields should be overwritten by *before* fields")
		}
		// there should be no information about duration and session end
		if stripped.RealtimeDuration != 0 ||
			stripped.LastRecordOfSession != false {
			t.Error("Stripped() returned record with too much information")
		}
	}
}

func TestGetCommandAndFirstWord(t *testing.T) {
	cmd, stWord, err := GetCommandAndFirstWord("cmd arg1 arg2")
	if err != nil || cmd != "cmd" || stWord != "cmd" {
		t.Error("GetCommandAndFirstWord() returned wrong Command OR FirstWord")
	}
}

func TestDistanceTo(t *testing.T) {
	paramsFull := DistParams{
		ExitCode:  1,
		MachineID: 1,
		SessionID: 1,
		Login:     1,
		Shell:     1,
		Pwd:       1,
		RealPwd:   1,
		Git:       1,
		Time:      1,
	}
	paramsZero := DistParams{}
	var prevRec EnrichedRecord
	for _, rec := range GetTestEnrichedRecords() {
		dist := rec.DistanceTo(rec, paramsFull)
		if dist != 0 {
			t.Error("DistanceTo() itself should be always 0")
		}
		dist = rec.DistanceTo(prevRec, paramsFull)
		if dist == 0 {
			t.Error("DistanceTo() between two test records shouldn't be 0")
		}
		dist = rec.DistanceTo(prevRec, paramsZero)
		if dist != 0 {
			t.Error("DistanceTo() should be 0 when DistParams is all zeros")
		}
		prevRec = rec
	}
}
