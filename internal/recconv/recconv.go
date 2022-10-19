package recconv

import (
	"fmt"

	"github.com/curusarn/resh/internal/record"
)

func LegacyToV1(r *record.Legacy) *record.V1 {
	return &record.V1{
		// FIXME: fill in all the fields

		// Flags: 0,

		CmdLine:  r.CmdLine,
		ExitCode: r.ExitCode,

		DeviceID:  r.ReshUUID,
		SessionID: r.SessionID,
		RecordID:  r.RecordID,

		Home:    r.Home,
		Pwd:     r.Pwd,
		RealPwd: r.RealPwd,

		// Logname:  r.Login,
		Device: r.Host,

		GitOriginRemote: r.GitOriginRemote,

		Time:     fmt.Sprintf("%.4f", r.RealtimeBefore),
		Duration: fmt.Sprintf("%.4f", r.RealtimeDuration),

		PartOne:        r.PartOne,
		PartsNotMerged: !r.PartsMerged,
	}
}
