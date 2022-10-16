package records

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/curusarn/resh/internal/histlist"
	"go.uber.org/zap"
)

// BaseRecord - common base for Record and FallbackRecord
type BaseRecord struct {
	// core
	CmdLine   string `json:"cmdLine"`
	ExitCode  int    `json:"exitCode"`
	Shell     string `json:"shell"`
	Uname     string `json:"uname"`
	SessionID string `json:"sessionId"`
	RecordID  string `json:"recordId"`

	// posix
	Home  string `json:"home"`
	Lang  string `json:"lang"`
	LcAll string `json:"lcAll"`
	Login string `json:"login"`
	//Path     string `json:"path"`
	Pwd      string `json:"pwd"`
	PwdAfter string `json:"pwdAfter"`
	ShellEnv string `json:"shellEnv"`
	Term     string `json:"term"`

	// non-posix"`
	RealPwd      string `json:"realPwd"`
	RealPwdAfter string `json:"realPwdAfter"`
	Pid          int    `json:"pid"`
	SessionPID   int    `json:"sessionPid"`
	Host         string `json:"host"`
	Hosttype     string `json:"hosttype"`
	Ostype       string `json:"ostype"`
	Machtype     string `json:"machtype"`
	Shlvl        int    `json:"shlvl"`

	// before after
	TimezoneBefore string `json:"timezoneBefore"`
	TimezoneAfter  string `json:"timezoneAfter"`

	RealtimeBefore      float64 `json:"realtimeBefore"`
	RealtimeAfter       float64 `json:"realtimeAfter"`
	RealtimeBeforeLocal float64 `json:"realtimeBeforeLocal"`
	RealtimeAfterLocal  float64 `json:"realtimeAfterLocal"`

	RealtimeDuration          float64 `json:"realtimeDuration"`
	RealtimeSinceSessionStart float64 `json:"realtimeSinceSessionStart"`
	RealtimeSinceBoot         float64 `json:"realtimeSinceBoot"`
	//Logs []string      `json: "logs"`

	GitDir               string `json:"gitDir"`
	GitRealDir           string `json:"gitRealDir"`
	GitOriginRemote      string `json:"gitOriginRemote"`
	GitDirAfter          string `json:"gitDirAfter"`
	GitRealDirAfter      string `json:"gitRealDirAfter"`
	GitOriginRemoteAfter string `json:"gitOriginRemoteAfter"`
	MachineID            string `json:"machineId"`

	OsReleaseID         string `json:"osReleaseId"`
	OsReleaseVersionID  string `json:"osReleaseVersionId"`
	OsReleaseIDLike     string `json:"osReleaseIdLike"`
	OsReleaseName       string `json:"osReleaseName"`
	OsReleasePrettyName string `json:"osReleasePrettyName"`

	ReshUUID     string `json:"reshUuid"`
	ReshVersion  string `json:"reshVersion"`
	ReshRevision string `json:"reshRevision"`

	// records come in two parts (collect and postcollect)
	PartOne     bool `json:"partOne,omitempty"` // false => part two
	PartsMerged bool `json:"partsMerged"`
	// special flag -> not an actual record but an session end
	SessionExit bool `json:"sessionExit,omitempty"`

	// recall metadata
	Recalled          bool     `json:"recalled"`
	RecallHistno      int      `json:"recallHistno,omitempty"`
	RecallStrategy    string   `json:"recallStrategy,omitempty"`
	RecallActionsRaw  string   `json:"recallActionsRaw,omitempty"`
	RecallActions     []string `json:"recallActions,omitempty"`
	RecallLastCmdLine string   `json:"recallLastCmdLine"`

	// recall command
	RecallPrefix string `json:"recallPrefix,omitempty"`

	// added by sanitizatizer
	Sanitized bool `json:"sanitized,omitempty"`
	CmdLength int  `json:"cmdLength,omitempty"`
}

// Record representing single executed command with its metadata
type Record struct {
	BaseRecord

	Cols  string `json:"cols"`
	Lines string `json:"lines"`
}

// EnrichedRecord - record enriched with additional data
type EnrichedRecord struct {
	Record

	// enriching fields - added "later"
	Command             string   `json:"command"`
	FirstWord           string   `json:"firstWord"`
	Invalid             bool     `json:"invalid"`
	SeqSessionID        uint64   `json:"seqSessionId"`
	LastRecordOfSession bool     `json:"lastRecordOfSession"`
	DebugThisRecord     bool     `json:"debugThisRecord"`
	Errors              []string `json:"errors"`
	// SeqSessionID uint64 `json:"seqSessionId,omitempty"`
}

// FallbackRecord when record is too old and can't be parsed into regular Record
type FallbackRecord struct {
	BaseRecord
	// older version of the record where cols and lines are int

	Cols  int `json:"cols"`  // notice the int type
	Lines int `json:"lines"` // notice the int type
}

// Convert from FallbackRecord to Record
func Convert(r *FallbackRecord) Record {
	return Record{
		BaseRecord: r.BaseRecord,
		// these two lines are the only reason we are doing this
		Cols:  strconv.Itoa(r.Cols),
		Lines: strconv.Itoa(r.Lines),
	}
}

// ToString - returns record the json
func (r EnrichedRecord) ToString() (string, error) {
	jsonRec, err := json.Marshal(r)
	if err != nil {
		return "marshalling error", err
	}
	return string(jsonRec), nil
}

// LoadFromFile loads records from 'fname' file
func LoadFromFile(sugar *zap.SugaredLogger, fname string) []Record {
	const allowedErrors = 3
	var encounteredErrors int
	var recs []Record
	file, err := os.Open(fname)
	if err != nil {
		sugar.Error("Failed to open resh history file - skipping reading resh history", zap.Error(err))
		return recs
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var i int
	for {
		var line string
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		}
		i++
		record := Record{}
		fallbackRecord := FallbackRecord{}
		err = json.Unmarshal([]byte(line), &record)
		if err != nil {
			err = json.Unmarshal([]byte(line), &fallbackRecord)
			if err != nil {
				encounteredErrors++
				sugar.Error("Could not decode line in resh history file",
					"lineContents", line,
					"lineNumber", i,
					zap.Error(err),
				)
				if encounteredErrors > allowedErrors {
					sugar.Fatal("Encountered too many errors during decoding - exiting",
						"allowedErrors", allowedErrors,
					)
				}
			}
			record = Convert(&fallbackRecord)
		}
		recs = append(recs, record)
	}
	if err != io.EOF {
		sugar.Error("Error while loading file", zap.Error(err))
	}
	sugar.Infow("Loaded resh history records",
		"recordCount", len(recs),
	)
	if encounteredErrors > 0 {
		// fix errors in the history file
		sugar.Warnw("Some history records could not be decoded - fixing resh history file by dropping them",
			"corruptedRecords", encounteredErrors,
		)
		fnameBak := fname + ".bak"
		sugar.Infow("Backing up current corrupted history file",
			"backupFilename", fnameBak,
		)
		err := copyFile(fname, fnameBak)
		if err != nil {
			sugar.Errorw("Failed to create a backup history file - aborting fixing history file",
				"backupFilename", fnameBak,
				zap.Error(err),
			)
			return recs
		}
		sugar.Info("Writing resh history file without errors ...")
		err = writeHistory(fname, recs)
		if err != nil {
			sugar.Errorw("Failed write fixed history file - aborting fixing history file",
				"filename", fname,
				zap.Error(err),
			)
		}
	}
	return recs
}

func copyFile(source, dest string) error {
	from, err := os.Open(source)
	if err != nil {
		return err
	}
	defer from.Close()

	// to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
	to, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
}

func writeHistory(fname string, history []Record) error {
	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, rec := range history {
		jsn, err := json.Marshal(rec)
		if err != nil {
			return fmt.Errorf("failed to encode record: %w", err)
		}
		file.Write(append(jsn, []byte("\n")...))
	}
	return nil
}

// LoadCmdLinesFromZshFile loads cmdlines from zsh history file
func LoadCmdLinesFromZshFile(sugar *zap.SugaredLogger, fname string) histlist.Histlist {
	hl := histlist.New(sugar)
	file, err := os.Open(fname)
	if err != nil {
		sugar.Error("Failed to open zsh history file - skipping reading zsh history", zap.Error(err))
		return hl
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// trim newline
		line = strings.TrimRight(line, "\n")
		var cmd string
		// zsh format EXTENDED_HISTORY
		// : 1576270617:0;make install
		// zsh format no EXTENDED_HISTORY
		// make install
		if len(line) == 0 {
			// skip empty
			continue
		}
		if strings.Contains(line, ":") && strings.Contains(line, ";") &&
			len(strings.Split(line, ":")) >= 3 && len(strings.Split(line, ";")) >= 2 {
			// contains at least 2x ':' and 1x ';' => assume EXTENDED_HISTORY
			cmd = strings.Split(line, ";")[1]
		} else {
			cmd = line
		}
		hl.AddCmdLine(cmd)
	}
	return hl
}

// LoadCmdLinesFromBashFile loads cmdlines from bash history file
func LoadCmdLinesFromBashFile(sugar *zap.SugaredLogger, fname string) histlist.Histlist {
	hl := histlist.New(sugar)
	file, err := os.Open(fname)
	if err != nil {
		sugar.Error("Failed to open bash history file - skipping reading bash history", zap.Error(err))
		return hl
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// trim newline
		line = strings.TrimRight(line, "\n")
		// trim spaces from left
		line = strings.TrimLeft(line, " ")
		// bash format (two lines)
		// #1576199174
		// make install
		if strings.HasPrefix(line, "#") {
			// is either timestamp or comment => skip
			continue
		}
		if len(line) == 0 {
			// skip empty
			continue
		}
		hl.AddCmdLine(line)
	}
	return hl
}
