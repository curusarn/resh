package strat

import (
	"github.com/curusarn/resh/pkg/records"
)

// ISimpleStrategy interface
type ISimpleStrategy interface {
	GetTitleAndDescription() (string, string)
	GetCandidates() []string
	AddHistoryRecord(record *records.EnrichedRecord) error
	ResetHistory() error
}

// IStrategy interface
type IStrategy interface {
	GetTitleAndDescription() (string, string)
	GetCandidates(r records.EnrichedRecord) []string
	AddHistoryRecord(record *records.EnrichedRecord) error
	ResetHistory() error
}

type simpleStrategyWrapper struct {
	strategy ISimpleStrategy
}

// NewSimpleStrategyWrapper returns IStrategy created by wrapping given ISimpleStrategy
func NewSimpleStrategyWrapper(strategy ISimpleStrategy) *simpleStrategyWrapper {
	return &simpleStrategyWrapper{strategy: strategy}
}

func (s *simpleStrategyWrapper) GetTitleAndDescription() (string, string) {
	return s.strategy.GetTitleAndDescription()
}

func (s *simpleStrategyWrapper) GetCandidates(r records.EnrichedRecord) []string {
	return s.strategy.GetCandidates()
}

func (s *simpleStrategyWrapper) AddHistoryRecord(r *records.EnrichedRecord) error {
	return s.strategy.AddHistoryRecord(r)
}

func (s *simpleStrategyWrapper) ResetHistory() error {
	return s.strategy.ResetHistory()
}
