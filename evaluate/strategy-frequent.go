package main

import (
	"sort"

	"github.com/curusarn/resh/common"
)

type strategyFrequent struct {
	history map[string]int
}

type strFrqEntry struct {
	cmdLine string
	count   int
}

func (s *strategyFrequent) init() {
	s.history = map[string]int{}
}

func (s *strategyFrequent) GetTitleAndDescription() (string, string) {
	return "frequent", "Use frequent commands"
}

func (s *strategyFrequent) GetCandidates() []string {
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

func (s *strategyFrequent) AddHistoryRecord(record *common.EnrichedRecord) error {
	s.history[record.CmdLine]++
	return nil
}

func (s *strategyFrequent) ResetHistory() error {
	s.history = map[string]int{}
	return nil
}
