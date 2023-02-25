package recutil

import (
	"errors"

	"github.com/curusarn/resh/internal/recordint"
	"github.com/curusarn/resh/record"
)

// TODO: reintroduce validation
// Validate returns error if the record is invalid
// func Validate(r *record.V1) error {
// 	if r.CmdLine == "" {
// 		return errors.New("There is no CmdLine")
// 	}
// 	if r.Time == 0 {
// 		return errors.New("There is no Time")
// 	}
// 	if r.RealPwd == "" {
// 		return errors.New("There is no Real Pwd")
// 	}
// 	if r.Pwd == "" {
// 		return errors.New("There is no Pwd")
// 	}
// 	return nil
// }

// TODO: maybe more to a more appropriate place
// TODO: cleanup the interface - stop modifying the part1 and returning a new record at the same time
// Merge two records (part1 - collect + part2 - postcollect)
func Merge(r1 *recordint.Collect, r2 *recordint.Collect) (record.V1, error) {
	if r1.SessionID != r2.SessionID {
		return record.V1{}, errors.New("Records to merge are not from the same session - r1:" + r1.SessionID + " r2:" + r2.SessionID)
	}
	if r1.Rec.RecordID != r2.Rec.RecordID {
		return record.V1{}, errors.New("Records to merge do not have the same ID - r1:" + r1.Rec.RecordID + " r2:" + r2.Rec.RecordID)
	}

	r := recordint.Collect{
		SessionID:  r1.SessionID,
		Shlvl:      r1.Shlvl,
		SessionPID: r1.SessionPID,

		Rec: r1.Rec,
	}
	r.Rec.ExitCode = r2.Rec.ExitCode
	r.Rec.Duration = r2.Rec.Duration
	r.Rec.PartOne = false
	r.Rec.PartsNotMerged = false
	return r.Rec, nil
}
