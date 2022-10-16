package recordint

// SearchApp record used for sending records to RESH-CLI
type SearchApp struct {
	IsRaw     bool
	SessionID string
	DeviceID  string

	CmdLine         string
	Host            string
	Pwd             string
	Home            string // helps us to collapse /home/user to tilde
	GitOriginRemote string
	ExitCode        int

	Time float64
}

// NewCliRecordFromCmdLine
func NewSearchAppFromCmdLine(cmdLine string) SearchApp {
	return SearchApp{
		IsRaw:   true,
		CmdLine: cmdLine,
	}
}

// NewCliRecord from EnrichedRecord
func NewSearchApp(r *Enriched) SearchApp {
	return SearchApp{
		IsRaw:           false,
		SessionID:       r.SessionID,
		CmdLine:         r.CmdLine,
		Host:            r.Hostname,
		Pwd:             r.Pwd,
		Home:            r.Home,
		GitOriginRemote: r.GitOriginRemote,
		ExitCode:        r.ExitCode,
		Time:            r.Time,
	}
}
