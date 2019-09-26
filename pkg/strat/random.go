package strat

import (
	"math/rand"
	"time"

	"github.com/curusarn/resh/pkg/records"
)

// Random prediction/recommendation strategy
type Random struct {
	CandidatesSize int
	history        []string
	historySet     map[string]bool
}

// Init see name
func (s *Random) Init() {
	s.history = nil
	s.historySet = map[string]bool{}
}

// GetTitleAndDescription see name
func (s *Random) GetTitleAndDescription() (string, string) {
	return "random", "Use random commands"
}

// GetCandidates see name
func (s *Random) GetCandidates() []string {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	var candidates []string
	candidateSet := map[string]bool{}
	for len(candidates) < s.CandidatesSize && len(candidates)*2 < len(s.historySet) {
		x := rand.Intn(len(s.history))
		candidate := s.history[x]
		if candidateSet[candidate] == false {
			candidateSet[candidate] = true
			candidates = append(candidates, candidate)
			continue
		}
	}
	return candidates
}

// AddHistoryRecord see name
func (s *Random) AddHistoryRecord(record *records.EnrichedRecord) error {
	s.history = append([]string{record.CmdLine}, s.history...)
	s.historySet[record.CmdLine] = true
	return nil
}

// ResetHistory see name
func (s *Random) ResetHistory() error {
	s.Init()
	return nil
}
