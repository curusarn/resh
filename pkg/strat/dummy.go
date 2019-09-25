package strat

import "github.com/curusarn/resh/pkg/records"

type Dummy struct {
	history []string
}

func (s *Dummy) GetTitleAndDescription() (string, string) {
	return "dummy", "Return empty candidate list"
}

func (s *Dummy) GetCandidates() []string {
	return nil
}

func (s *Dummy) AddHistoryRecord(record *records.EnrichedRecord) error {
	s.history = append(s.history, record.CmdLine)
	return nil
}

func (s *Dummy) ResetHistory() error {
	return nil
}
