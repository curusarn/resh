package strat

import "github.com/curusarn/resh/pkg/records"

// Dummy prediction/recommendation strategy
type Dummy struct {
	history []string
}

// GetTitleAndDescription see name
func (s *Dummy) GetTitleAndDescription() (string, string) {
	return "dummy", "Return empty candidate list"
}

// GetCandidates see name
func (s *Dummy) GetCandidates() []string {
	return nil
}

// AddHistoryRecord see name
func (s *Dummy) AddHistoryRecord(record *records.EnrichedRecord) error {
	s.history = append(s.history, record.CmdLine)
	return nil
}

// ResetHistory see name
func (s *Dummy) ResetHistory() error {
	return nil
}
