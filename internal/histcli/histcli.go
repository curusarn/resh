package histcli

import (
	"github.com/curusarn/resh/internal/recordint"
	"sync"
)

// Histcli is a dump of history preprocessed for resh cli purposes
type Histcli struct {
	// list of records
	list []recordint.SearchApp
	lock sync.RWMutex
}

// New Histcli
func New() *Histcli {
	return &Histcli{}
}

// AddRecord to the histcli
func (h *Histcli) AddRecord(rec *recordint.Indexed) {
	cli := recordint.NewSearchApp(rec)
	h.lock.Lock()
	defer h.lock.Unlock()

	h.list = append(h.list, cli)
}

// AddCmdLine to the histcli
func (h *Histcli) AddCmdLine(cmdline string) {
	cli := recordint.NewSearchAppFromCmdLine(cmdline)
	h.lock.Lock()
	defer h.lock.Unlock()

	h.list = append(h.list, cli)
}

func (h *Histcli) Dump() []recordint.SearchApp {
	h.lock.RLock()
	defer h.lock.RUnlock()

	return h.list
}
