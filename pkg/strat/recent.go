package strat

import "github.com/curusarn/resh/pkg/records"

// Recent prediction/recommendation strategy
type Recent struct {
	history []string
}

// GetTitleAndDescription see name
func (s *Recent) GetTitleAndDescription() (string, string) {
	return "recent", "Use recent commands"
}

// GetCandidates see name
func (s *Recent) GetCandidates() []string {
	return s.history
}

// AddHistoryRecord see name
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

// ResetHistory see name
func (s *Recent) ResetHistory() error {
	s.history = nil
	return nil
}
