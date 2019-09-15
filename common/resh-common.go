package common

import (
	"log"
	"strconv"
	"strings"

	"github.com/mattn/go-shellwords"
)

// BaseRecord - common base for Record and FallbackRecord
type BaseRecord struct {
	// core
	CmdLine   string `json:"cmdLine"`
	ExitCode  int    `json:"exitCode"`
	Shell     string `json:"shell"`
	Uname     string `json:"uname"`
	SessionID string `json:"sessionId"`

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
	SessionPid   int    `json:"sessionPid"`
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

	GitDir          string `json:"gitDir"`
	GitRealDir      string `json:"gitRealDir"`
	GitOriginRemote string `json:"gitOriginRemote"`
	MachineID       string `json:"machineId"`

	OsReleaseID         string `json:"osReleaseId"`
	OsReleaseVersionID  string `json:"osReleaseVersionId"`
	OsReleaseIDLike     string `json:"osReleaseIdLike"`
	OsReleaseName       string `json:"osReleaseName"`
	OsReleasePrettyName string `json:"osReleasePrettyName"`

	ReshUUID     string `json:"reshUuid"`
	ReshVersion  string `json:"reshVersion"`
	ReshRevision string `json:"reshRevision"`

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
	Command      string `json:"command"`
	FirstWord    string `json:"firstWord"`
	Invalid      bool   `json:"invalid"`
	SeqSessionID uint64 `json:"seqSessionId"`
	// SeqSessionID uint64 `json:"seqSessionId,omitempty"`
}

// FallbackRecord when record is too old and can't be parsed into regular Record
type FallbackRecord struct {
	BaseRecord
	// older version of the record where cols and lines are int

	Cols  int `json:"cols"`  // notice the int type
	Lines int `json:"lines"` // notice the int type
}

// ConvertRecord from FallbackRecord to Record
func ConvertRecord(r *FallbackRecord) Record {
	return Record{
		BaseRecord: r.BaseRecord,
		// these two lines are the only reason we are doing this
		Cols:  strconv.Itoa(r.Cols),
		Lines: strconv.Itoa(r.Lines),
	}
}

// Enrich - adds additional fields to the record
func (r Record) Enrich() EnrichedRecord {
	record := EnrichedRecord{Record: r}
	// Get command/first word from commandline
	record.Command, record.FirstWord = GetCommandAndFirstWord(r.CmdLine)
	err := r.Validate()
	if err != nil {
		log.Println("Invalid command:", r.CmdLine)
		record.Invalid = true
	}
	return record
	// TODO: Detect and mark simple commands r.Simple
}

// Validate - returns error if the record is invalid
func (r *Record) Validate() error {
	return nil
}

// GetCommandAndFirstWord func
func GetCommandAndFirstWord(cmdLine string) (string, string) {
	args, err := shellwords.Parse(cmdLine)
	if err != nil {
		log.Println("shellwords Error:", err, " (cmdLine: <", cmdLine, "> )")
		return "<shellwords_error>", "<shellwords_error>"
	}
	if len(args) == 0 {
		return "", ""
	}
	i := 0
	for true {
		// commands in shell sometimes look like this `variable=something command argument otherArgument --option`
		//		to get the command we skip over tokens that contain '='
		if strings.ContainsRune(args[i], '=') && len(args) > i+1 {
			i++
			continue
		}
		return args[i], args[0]
	}
	log.Fatal("GetCommandAndFirstWord error: this should not happen!")
	return "ERROR", "ERROR"
}

// Config struct
type Config struct {
	Port int
}
