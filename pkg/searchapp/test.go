package searchapp

import (
	"github.com/curusarn/resh/pkg/histcli"
	"github.com/curusarn/resh/pkg/msg"
	"github.com/curusarn/resh/pkg/records"
)

// LoadHistoryFromFile ...
func LoadHistoryFromFile(historyPath string, numLines int) msg.CliResponse {
	recs := records.LoadFromFile(historyPath)
	if numLines != 0 && numLines < len(recs) {
		recs = recs[:numLines]
	}
	cliRecords := histcli.New()
	for i := len(recs) - 1; i >= 0; i-- {
		rec := recs[i]
		cliRecords.AddRecord(rec)
	}
	return msg.CliResponse{CliRecords: cliRecords.List}
}
