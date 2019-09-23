package main

import (
	"sort"
	"strconv"

	"github.com/curusarn/resh/common"
)

type strategyRecordDistance struct {
	history    []common.EnrichedRecord
	distParams common.DistParams
	maxDepth   int
	label      string
}

type strDistEntry struct {
	cmdLine  string
	distance float64
}

func (s *strategyRecordDistance) init() {
	s.history = nil
}

func (s *strategyRecordDistance) GetTitleAndDescription() (string, string) {
	return "record distance (depth:" + strconv.Itoa(s.maxDepth) + ";" + s.label + ")", "Use record distance to recommend commands"
}

func (s *strategyRecordDistance) GetCandidates() []string {
	if len(s.history) == 0 {
		return nil
	}
	var prevRecord common.EnrichedRecord
	prevRecord = s.history[0]
	prevRecord.SetCmdLine("")
	prevRecord.SetBeforeToAfter()
	var mapItems []strDistEntry
	for i, record := range s.history {
		if s.maxDepth != 0 && i > s.maxDepth {
			break
		}
		distance := record.DistanceTo(prevRecord, s.distParams)
		mapItems = append(mapItems, strDistEntry{record.CmdLine, distance})
	}
	sort.SliceStable(mapItems, func(i int, j int) bool { return mapItems[i].distance < mapItems[j].distance })
	var hist []string
	histSet := map[string]bool{}
	for _, item := range mapItems {
		if histSet[item.cmdLine] {
			continue
		}
		histSet[item.cmdLine] = true
		hist = append(hist, item.cmdLine)
	}
	return hist
}

func (s *strategyRecordDistance) AddHistoryRecord(record *common.EnrichedRecord) error {
	// append record to front
	s.history = append([]common.EnrichedRecord{*record}, s.history...)
	return nil
}

func (s *strategyRecordDistance) ResetHistory() error {
	s.init()
	return nil
}
