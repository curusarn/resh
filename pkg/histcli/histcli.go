package histcli

import (
	"github.com/curusarn/resh/pkg/records"
)

// Histcli is a dump of history preprocessed for resh cli purposes
type Histcli struct {
	// list of records
	List []records.CliRecord
}

// New Histcli
func New() Histcli {
	return Histcli{}
}

// AddRecord to the histcli
func (h *Histcli) AddRecord(record records.Record) {
	enriched := records.Enriched(record)
	cli := records.NewCliRecord(enriched)

	h.List = append(h.List, cli)
}

// AddCmdLine to the histcli
func (h *Histcli) AddCmdLine(cmdline string) {
	cli := records.NewCliRecordFromCmdLine(cmdline)

	h.List = append(h.List, cli)
}
