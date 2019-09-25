package strat

import (
	"sort"

	"github.com/curusarn/resh/pkg/records"
)

type Frequent struct {
	history map[string]int
}

type strFrqEntry struct {
	cmdLine string
	count   int
}

func (s *Frequent) init() {
	s.history = map[string]int{}
}

func (s *Frequent) GetTitleAndDescription() (string, string) {
	return "frequent", "Use frequent commands"
}

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

func (s *Frequent) AddHistoryRecord(record *records.EnrichedRecord) error {
	s.history[record.CmdLine]++
	return nil
}

func (s *Frequent) ResetHistory() error {
	s.init()
	return nil
}
