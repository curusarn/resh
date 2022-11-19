package histcli

import (
	"github.com/curusarn/resh/internal/recordint"
	"go.uber.org/zap"
	"sync"
)

// Histcli is a dump of history preprocessed for resh cli purposes
type Histcli struct {
	// list of records
	list     []recordint.SearchApp
	knownIds map[string]struct{}
	lock     sync.RWMutex
	sugar    *zap.SugaredLogger
	latest   map[string]float64
}

// New Histcli
func New(sugar *zap.SugaredLogger) *Histcli {
	return &Histcli{
		sugar:    sugar.With(zap.String("component", "histCli")),
		knownIds: map[string]struct{}{},
		latest:   map[string]float64{},
	}
}

// AddRecord to the histcli
func (h *Histcli) AddRecord(rec *recordint.Indexed) {
	cli := recordint.NewSearchApp(rec)
	h.lock.Lock()
	defer h.lock.Unlock()

	if _, ok := h.knownIds[rec.Rec.RecordID]; !ok {
		h.knownIds[rec.Rec.RecordID] = struct{}{}
		h.list = append(h.list, cli)
		h.updateLatestPerDevice(cli)
	} else {
		h.sugar.Debugw("Record is already present", "id", rec.Rec.RecordID)
	}
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

// updateLatestPerDevice should be called only with write lock because it does not lock on its own.
func (h *Histcli) updateLatestPerDevice(rec recordint.SearchApp) {
	if l, ok := h.latest[rec.DeviceID]; ok {
		if rec.Time > l {
			h.latest[rec.DeviceID] = rec.Time
		}
	} else {
		h.latest[rec.DeviceID] = rec.Time
	}
}

func (h *Histcli) LatestRecordsPerDevice() map[string]float64 {
	h.lock.RLock()
	defer h.lock.RUnlock()

	return h.latest
}
