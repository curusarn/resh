package main

import "github.com/curusarn/resh/pkg/records"

type strategyDummy struct {
	history []string
}

func (s *strategyDummy) GetTitleAndDescription() (string, string) {
	return "dummy", "Return empty candidate list"
}

func (s *strategyDummy) GetCandidates() []string {
	return nil
}

func (s *strategyDummy) AddHistoryRecord(record *records.EnrichedRecord) error {
	s.history = append(s.history, record.CmdLine)
	return nil
}

func (s *strategyDummy) ResetHistory() error {
	return nil
}