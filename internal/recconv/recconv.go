package recconv

import (
	"github.com/curusarn/resh/internal/record"
)

func LegacyToV1(r *record.Legacy) *record.V1 {
	return &record.V1{
		// FIXME: fill in all the fields

		// Flags: 0,

		DeviceID:  r.MachineID,
		SessionID: r.SessionID,
		RecordID:  r.RecordID,

		CmdLine:  r.CmdLine,
		ExitCode: r.ExitCode,

		Home:    r.Home,
		Pwd:     r.Pwd,
		RealPwd: r.RealPwd,

		Logname:  r.Login,
		Hostname: r.Host,

		GitOriginRemote: r.GitOriginRemote,

		Time:     r.RealtimeBefore,
		Duration: r.RealtimeDuration,

		PartOne:        r.PartOne,
		PartsNotMerged: !r.PartsMerged,
	}
}
