package searchapp

import (
	"github.com/curusarn/resh/internal/histcli"
	"github.com/curusarn/resh/internal/msg"
	"github.com/curusarn/resh/internal/records"
	"go.uber.org/zap"
)

// LoadHistoryFromFile ...
func LoadHistoryFromFile(sugar *zap.SugaredLogger, historyPath string, numLines int) msg.CliResponse {
	recs := records.LoadFromFile(sugar, historyPath)
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
