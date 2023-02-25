package recordint

import (
	"strconv"

	"github.com/curusarn/resh/internal/normalize"
	"github.com/curusarn/resh/record"
	"go.uber.org/zap"
)

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

	// file index
	Idx int
}

func NewSearchAppFromCmdLine(cmdLine string) SearchApp {
	return SearchApp{
		IsRaw:   true,
		CmdLine: cmdLine,
	}
}

// The error handling here could be better
func NewSearchApp(sugar *zap.SugaredLogger, r *record.V1) SearchApp {
	time, err := strconv.ParseFloat(r.Time, 64)
	if err != nil {
		sugar.Errorw("Error while parsing time as float", zap.Error(err),
			"time", time)
	}
	return SearchApp{
		IsRaw:     false,
		SessionID: r.SessionID,
		CmdLine:   r.CmdLine,
		Host:      r.Device,
		Pwd:       r.Pwd,
		Home:      r.Home,
		// TODO: is this the right place to normalize the git remote?
		GitOriginRemote: normalize.GitRemote(sugar, r.GitOriginRemote),
		ExitCode:        r.ExitCode,
		Time:            time,
	}
}
