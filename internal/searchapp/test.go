package searchapp

import (
	"github.com/curusarn/resh/internal/histcli"
	"github.com/curusarn/resh/internal/msg"
	"github.com/curusarn/resh/internal/recio"
	"go.uber.org/zap"
)

// LoadHistoryFromFile ...
func LoadHistoryFromFile(sugar *zap.SugaredLogger, historyPath string, numLines int) msg.CliResponse {
	rio := recio.New(sugar)
	recs, _, err := rio.ReadFile(historyPath)
	if err != nil {
		sugar.Panicf("failed to read hisotry file: %w", err)
	}
	if numLines != 0 && numLines < len(recs) {
		recs = recs[:numLines]
	}
	cliRecords := histcli.New(sugar)
	for i := len(recs) - 1; i >= 0; i-- {
		rec := recs[i]
		cliRecords.AddRecord(&rec)
	}
	return msg.CliResponse{Records: cliRecords.Dump()}
}
