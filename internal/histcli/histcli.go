package histcli

import (
	"github.com/curusarn/resh/internal/recordint"
)

// Histcli is a dump of history preprocessed for resh cli purposes
type Histcli struct {
	// list of records
	List []recordint.SearchApp
}

// New Histcli
func New() Histcli {
	return Histcli{}
}

// AddRecord to the histcli
func (h *Histcli) AddRecord(rec *recordint.Indexed) {
	cli := recordint.NewSearchApp(rec)

	h.List = append(h.List, cli)
}

// AddCmdLine to the histcli
func (h *Histcli) AddCmdLine(cmdline string) {
	cli := recordint.NewSearchAppFromCmdLine(cmdline)

	h.List = append(h.List, cli)
}
