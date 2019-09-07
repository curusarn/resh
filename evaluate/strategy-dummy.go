package main

import "github.com/curusarn/resh/common"

type strategyDummy struct {
	history []string
}

func (s *strategyDummy) GetTitleAndDescription() (string, string) {
	return "recent", "Use recent commands"
}

func (s *strategyDummy) GetCandidates() []string {
	return nil
}

func (s *strategyDummy) AddHistoryRecord(record *common.Record) error {
	s.history = append(s.history, record.CmdLine)
	return nil
}

func (s *strategyDummy) ResetHistory() error {
	return nil
}
