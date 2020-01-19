package histcli

import (
	"github.com/curusarn/resh/pkg/records"
)

// Histcli is a dump of history preprocessed for resh cli purposes
type Histcli struct {
	// list of records
	List []records.EnrichedRecord
}

// New Histcli
func New() Histcli {
	return Histcli{}
}

// AddRecord to the histcli
func (h *Histcli) AddRecord(record records.Record) {
	enriched := records.Enriched(record)

	h.List = append(h.List, enriched)
}
