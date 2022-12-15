package histcli

import (
	"github.com/curusarn/resh/internal/recordint"
	"github.com/curusarn/resh/record"
	"go.uber.org/zap"
)

// Histcli is a dump of history preprocessed for resh cli purposes
type Histcli struct {
	// list of records
	List []recordint.SearchApp

	sugar *zap.SugaredLogger
}

// New Histcli
func New(sugar *zap.SugaredLogger) Histcli {
	return Histcli{}
}

// AddRecord to the histcli
func (h *Histcli) AddRecord(rec *record.V1) {
	cli := recordint.NewSearchApp(h.sugar, rec)

	h.List = append(h.List, cli)
}

// AddCmdLine to the histcli
func (h *Histcli) AddCmdLine(cmdline string) {
	cli := recordint.NewSearchAppFromCmdLine(cmdline)

	h.List = append(h.List, cli)
}
