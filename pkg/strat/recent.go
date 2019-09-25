package strat

import "github.com/curusarn/resh/pkg/records"

type Recent struct {
	history []string
}

func (s *Recent) GetTitleAndDescription() (string, string) {
	return "recent", "Use recent commands"
}

func (s *Recent) GetCandidates() []string {
	return s.history
}

func (s *Recent) AddHistoryRecord(record *records.EnrichedRecord) error {
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

func (s *Recent) ResetHistory() error {
	s.history = nil
	return nil
}
