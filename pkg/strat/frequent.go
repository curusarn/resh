package strat

import (
	"sort"

	"github.com/curusarn/resh/pkg/records"
)

// Frequent prediction/recommendation strategy
type Frequent struct {
	history map[string]int
}

type strFrqEntry struct {
	cmdLine string
	count   int
}

// Init see name
func (s *Frequent) Init() {
	s.history = map[string]int{}
}

// GetTitleAndDescription see name
func (s *Frequent) GetTitleAndDescription() (string, string) {
	return "frequent", "Use frequent commands"
}

// GetCandidates see name
func (s *Frequent) GetCandidates() []string {
	var mapItems []strFrqEntry
	for cmdLine, count := range s.history {
		mapItems = append(mapItems, strFrqEntry{cmdLine, count})
	}
	sort.Slice(mapItems, func(i int, j int) bool { return mapItems[i].count > mapItems[j].count })
	var hist []string
	for _, item := range mapItems {
		hist = append(hist, item.cmdLine)
	}
	return hist
}

// AddHistoryRecord see name
func (s *Frequent) AddHistoryRecord(record *records.EnrichedRecord) error {
	s.history[record.CmdLine]++
	return nil
}

// ResetHistory see name
func (s *Frequent) ResetHistory() error {
	s.Init()
	return nil
}
