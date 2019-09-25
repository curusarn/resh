package strat

import (
	"math/rand"
	"time"

	"github.com/curusarn/resh/pkg/records"
)

type Random struct {
	CandidatesSize int
	history        []string
	historySet     map[string]bool
}

func (s *Random) Init() {
	s.history = nil
	s.historySet = map[string]bool{}
}

func (s *Random) GetTitleAndDescription() (string, string) {
	return "random", "Use random commands"
}

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

func (s *Random) AddHistoryRecord(record *records.EnrichedRecord) error {
	s.history = append([]string{record.CmdLine}, s.history...)
	s.historySet[record.CmdLine] = true
	return nil
}

func (s *Random) ResetHistory() error {
	s.Init()
	return nil
}
