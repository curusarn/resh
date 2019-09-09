package common

import (
	"log"
	"strconv"

	"github.com/mattn/go-shellwords"
)

// Record representing single executed command with its metadata
type Record struct {
	// core
	CmdLine   string `json:"cmdLine"`
	ExitCode  int    `json:"exitCode"`
	Shell     string `json:"shell"`
	Uname     string `json:"uname"`
	SessionId string `json:"sessionId"`

	// posix
	Cols  string `json:"cols"`
	Lines string `json:"lines"`
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
	MachineId       string `json:"machineId"`

	OsReleaseId         string `json:"osReleaseId"`
	OsReleaseVersionId  string `json:"osReleaseVersionId"`
	OsReleaseIdLike     string `json:"osReleaseIdLike"`
	OsReleaseName       string `json:"osReleaseName"`
	OsReleasePrettyName string `json:"osReleasePrettyName"`

	ReshUuid     string `json:"reshUuid"`
	ReshVersion  string `json:"reshVersion"`
	ReshRevision string `json:"reshRevision"`

	// added by sanitizatizer
	Sanitized bool `json:"sanitized"`
	CmdLength int  `json:"cmdLength"`

	// enriching fields - added "later"
	FirstWord string `json:"firstWord"`
}

// FallbackRecord when record is too old and can't be parsed into regular Record
type FallbackRecord struct {
	// older version of the record where cols and lines are int

	// core
	CmdLine   string `json:"cmdLine"`
	ExitCode  int    `json:"exitCode"`
	Shell     string `json:"shell"`
	Uname     string `json:"uname"`
	SessionId string `json:"sessionId"`

	// posix
	Cols  int    `json:"cols"`  // notice the in type
	Lines int    `json:"lines"` // notice the in type
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
	MachineId       string `json:"machineId"`

	OsReleaseId         string `json:"osReleaseId"`
	OsReleaseVersionId  string `json:"osReleaseVersionId"`
	OsReleaseIdLike     string `json:"osReleaseIdLike"`
	OsReleaseName       string `json:"osReleaseName"`
	OsReleasePrettyName string `json:"osReleasePrettyName"`

	ReshUuid     string `json:"reshUuid"`
	ReshVersion  string `json:"reshVersion"`
	ReshRevision string `json:"reshRevision"`
}

// ConvertRecord from FallbackRecord to Record
func ConvertRecord(r *FallbackRecord) Record {
	return Record{
		// core
		CmdLine:   r.CmdLine,
		ExitCode:  r.ExitCode,
		Shell:     r.Shell,
		Uname:     r.Uname,
		SessionId: r.SessionId,

		// posix
		// these two lines are the only reason we are doing this
		Cols:  strconv.Itoa(r.Cols),
		Lines: strconv.Itoa(r.Lines),

		Home:  r.Home,
		Lang:  r.Lang,
		LcAll: r.LcAll,
		Login: r.Login,
		// Path:     r.path,
		Pwd:      r.Pwd,
		PwdAfter: r.PwdAfter,
		ShellEnv: r.ShellEnv,
		Term:     r.Term,

		// non-posix
		RealPwd:      r.RealPwd,
		RealPwdAfter: r.RealPwdAfter,
		Pid:          r.Pid,
		SessionPid:   r.SessionPid,
		Host:         r.Host,
		Hosttype:     r.Hosttype,
		Ostype:       r.Ostype,
		Machtype:     r.Machtype,
		Shlvl:        r.Shlvl,

		// before after
		TimezoneBefore: r.TimezoneBefore,
		TimezoneAfter:  r.TimezoneAfter,

		RealtimeBefore:      r.RealtimeBefore,
		RealtimeAfter:       r.RealtimeAfter,
		RealtimeBeforeLocal: r.RealtimeBeforeLocal,
		RealtimeAfterLocal:  r.RealtimeAfterLocal,

		RealtimeDuration:          r.RealtimeDuration,
		RealtimeSinceSessionStart: r.RealtimeSinceSessionStart,
		RealtimeSinceBoot:         r.RealtimeSinceBoot,

		GitDir:          r.GitDir,
		GitRealDir:      r.GitRealDir,
		GitOriginRemote: r.GitOriginRemote,
		MachineId:       r.MachineId,

		OsReleaseId:         r.OsReleaseId,
		OsReleaseVersionId:  r.OsReleaseVersionId,
		OsReleaseIdLike:     r.OsReleaseIdLike,
		OsReleaseName:       r.OsReleaseName,
		OsReleasePrettyName: r.OsReleasePrettyName,

		ReshUuid:     r.ReshUuid,
		ReshVersion:  r.ReshVersion,
		ReshRevision: r.ReshRevision,
	}
}

// Enrich - adds additional fields to the record
func (r *Record) Enrich() {
	// Get command/first word from commandline
	r.FirstWord = GetCommandFromCommandLine(r.CmdLine)
}

// GetCommandFromCommandLine func
func GetCommandFromCommandLine(cmdLine string) string {
	args, err := shellwords.Parse(cmdLine)
	if err != nil {
		log.Fatal("shellwords Error:", err)
	}
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

// Config struct
type Config struct {
	Port int
}
