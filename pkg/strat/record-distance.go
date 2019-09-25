package strat

import (
	"sort"
	"strconv"

	"github.com/curusarn/resh/pkg/records"
)

type RecordDistance struct {
	history    []records.EnrichedRecord
	DistParams records.DistParams
	MaxDepth   int
	Label      string
}

type strDistEntry struct {
	cmdLine  string
	distance float64
}

func (s *RecordDistance) Init() {
	s.history = nil
}

func (s *RecordDistance) GetTitleAndDescription() (string, string) {
	return "record distance (depth:" + strconv.Itoa(s.MaxDepth) + ";" + s.Label + ")", "Use record distance to recommend commands"
}

func (s *RecordDistance) GetCandidates(strippedRecord records.EnrichedRecord) []string {
	if len(s.history) == 0 {
		return nil
	}
	var mapItems []strDistEntry
	for i, record := range s.history {
		if s.MaxDepth != 0 && i > s.MaxDepth {
			break
		}
		distance := record.DistanceTo(strippedRecord, s.DistParams)
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

func (s *RecordDistance) AddHistoryRecord(record *records.EnrichedRecord) error {
	// append record to front
	s.history = append([]records.EnrichedRecord{*record}, s.history...)
	return nil
}

func (s *RecordDistance) ResetHistory() error {
	s.Init()
	return nil
}
