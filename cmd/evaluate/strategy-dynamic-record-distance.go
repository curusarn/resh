package main

import (
	"math"
	"sort"
	"strconv"

	"github.com/curusarn/resh/pkg/records"
)

type strategyDynamicRecordDistance struct {
	history            []records.EnrichedRecord
	distParams         records.DistParams
	pwdHistogram       map[string]int
	realPwdHistogram   map[string]int
	gitOriginHistogram map[string]int
	maxDepth           int
	label              string
}

type strDynDistEntry struct {
	cmdLine  string
	distance float64
}

func (s *strategyDynamicRecordDistance) init() {
	s.history = nil
	s.pwdHistogram = map[string]int{}
	s.realPwdHistogram = map[string]int{}
	s.gitOriginHistogram = map[string]int{}
}

func (s *strategyDynamicRecordDistance) GetTitleAndDescription() (string, string) {
	return "dynamic record distance (depth:" + strconv.Itoa(s.maxDepth) + ";" + s.label + ")", "Use TF-IDF record distance to recommend commands"
}

func (s *strategyDynamicRecordDistance) idf(count int) float64 {
	return math.Log(float64(len(s.history)) / float64(count))
}

func (s *strategyDynamicRecordDistance) GetCandidates(strippedRecord records.EnrichedRecord) []string {
	if len(s.history) == 0 {
		return nil
	}
	var mapItems []strDynDistEntry
	for i, record := range s.history {
		if s.maxDepth != 0 && i > s.maxDepth {
			break
		}
		distParams := records.DistParams{
			Pwd:       s.distParams.Pwd * s.idf(s.pwdHistogram[strippedRecord.PwdAfter]),
			RealPwd:   s.distParams.RealPwd * s.idf(s.realPwdHistogram[strippedRecord.RealPwdAfter]),
			Git:       s.distParams.Git * s.idf(s.gitOriginHistogram[strippedRecord.GitOriginRemote]),
			Time:      s.distParams.Time,
			SessionID: s.distParams.SessionID,
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

func (s *strategyDynamicRecordDistance) AddHistoryRecord(record *records.EnrichedRecord) error {
	// append record to front
	s.history = append([]records.EnrichedRecord{*record}, s.history...)
	s.pwdHistogram[record.Pwd]++
	s.realPwdHistogram[record.RealPwd]++
	s.gitOriginHistogram[record.GitOriginRemote]++
	return nil
}

func (s *strategyDynamicRecordDistance) ResetHistory() error {
	s.init()
	return nil
}
