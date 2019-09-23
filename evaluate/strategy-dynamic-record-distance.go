package main

import (
	"math"
	"sort"
	"strconv"

	"github.com/curusarn/resh/common"
)

type strategyDynamicRecordDistance struct {
	history          []common.EnrichedRecord
	distParams       common.DistParams
	pwdHistogram     map[string]int
	realPwdHistogram map[string]int
	maxDepth         int
	label            string
}

type strDynDistEntry struct {
	cmdLine  string
	distance float64
}

func (s *strategyDynamicRecordDistance) init() {
	s.history = nil
	s.pwdHistogram = map[string]int{}
	s.realPwdHistogram = map[string]int{}
}

func (s *strategyDynamicRecordDistance) GetTitleAndDescription() (string, string) {
	return "dynamic record distance (depth:" + strconv.Itoa(s.maxDepth) + ";" + s.label + ")", "Use TF-IDF record distance to recommend commands"
}

func (s *strategyDynamicRecordDistance) idf(count int) float64 {
	return math.Log(float64(len(s.history)) / float64(count))
}

func (s *strategyDynamicRecordDistance) GetCandidates() []string {
	if len(s.history) == 0 {
		return nil
	}
	var prevRecord common.EnrichedRecord
	prevRecord = s.history[0]
	prevRecord.SetCmdLine("")
	prevRecord.SetBeforeToAfter()
	var mapItems []strDynDistEntry
	for i, record := range s.history {
		if s.maxDepth != 0 && i > s.maxDepth {
			break
		}
		distParams := common.DistParams{
			Pwd:       s.distParams.Pwd * s.idf(s.pwdHistogram[prevRecord.PwdAfter]),
			RealPwd:   s.distParams.RealPwd * s.idf(s.realPwdHistogram[prevRecord.RealPwdAfter]),
			Time:      s.distParams.Time,
			SessionID: s.distParams.SessionID,
		}
		distance := record.DistanceTo(prevRecord, distParams)
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

func (s *strategyDynamicRecordDistance) AddHistoryRecord(record *common.EnrichedRecord) error {
	// append record to front
	s.history = append([]common.EnrichedRecord{*record}, s.history...)
	s.pwdHistogram[record.Pwd]++
	s.realPwdHistogram[record.RealPwd]++
	return nil
}

func (s *strategyDynamicRecordDistance) ResetHistory() error {
	s.init()
	return nil
}
