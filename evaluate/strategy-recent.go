package main

import "github.com/curusarn/resh/common"

type strategyRecent struct {
	history []string
}

func (s *strategyRecent) GetTitleAndDescription() (string, string) {
	return "recent", "Use recent commands"
}

func (s *strategyRecent) GetCandidates() []string {
	return s.history
}

func (s *strategyRecent) AddHistoryRecord(record *common.Record) error {
	// remove previous occurance of record
	for i, cmd := range s.history {
		if cmd == record.CmdLine {
			s.history = append(s.history[:i], s.history[i+1:]...)
		}
	}
	// append new record
	s.history = append([]string{record.CmdLine}, s.history...)
	return nil
}

func (s *strategyRecent) ResetHistory() error {
	s.history = nil
	return nil
}