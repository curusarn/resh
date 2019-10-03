package strat

import (
	"math"
	"sort"
	"strconv"

	"github.com/curusarn/resh/pkg/records"
)

// DynamicRecordDistance prediction/recommendation strategy
type DynamicRecordDistance struct {
	history            []records.EnrichedRecord
	DistParams         records.DistParams
	pwdHistogram       map[string]int
	realPwdHistogram   map[string]int
	gitOriginHistogram map[string]int
	MaxDepth           int
	Label              string
}

type strDynDistEntry struct {
	cmdLine  string
	distance float64
}

// Init see name
func (s *DynamicRecordDistance) Init() {
	s.history = nil
	s.pwdHistogram = map[string]int{}
	s.realPwdHistogram = map[string]int{}
	s.gitOriginHistogram = map[string]int{}
}

// GetTitleAndDescription see name
func (s *DynamicRecordDistance) GetTitleAndDescription() (string, string) {
	return "dynamic record distance (depth:" + strconv.Itoa(s.MaxDepth) + ";" + s.Label + ")", "Use TF-IDF record distance to recommend commands"
}

func (s *DynamicRecordDistance) idf(count int) float64 {
	return math.Log(float64(len(s.history)) / float64(count))
}

// GetCandidates see name
func (s *DynamicRecordDistance) GetCandidates(strippedRecord records.EnrichedRecord) []string {
	if len(s.history) == 0 {
		return nil
	}
	var mapItems []strDynDistEntry
	for i, record := range s.history {
		if s.MaxDepth != 0 && i > s.MaxDepth {
			break
		}
		distParams := records.DistParams{
			Pwd:       s.DistParams.Pwd * s.idf(s.pwdHistogram[strippedRecord.PwdAfter]),
			RealPwd:   s.DistParams.RealPwd * s.idf(s.realPwdHistogram[strippedRecord.RealPwdAfter]),
			Git:       s.DistParams.Git * s.idf(s.gitOriginHistogram[strippedRecord.GitOriginRemote]),
			Time:      s.DistParams.Time,
			SessionID: s.DistParams.SessionID,
		}
		distance := record.DistanceTo(strippedRecord, distParams)
		mapItems = append(mapItems, strDynDistEntry{record.CmdLine, distance})
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

// AddHistoryRecord see name
func (s *DynamicRecordDistance) AddHistoryRecord(record *records.EnrichedRecord) error {
	// append record to front
	s.history = append([]records.EnrichedRecord{*record}, s.history...)
	s.pwdHistogram[record.Pwd]++
	s.realPwdHistogram[record.RealPwd]++
	s.gitOriginHistogram[record.GitOriginRemote]++
	return nil
}

// ResetHistory see name
func (s *DynamicRecordDistance) ResetHistory() error {
	s.Init()
	return nil
}
