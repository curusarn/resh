package main

import (
	"math/rand"
	"time"

	"github.com/curusarn/resh/pkg/records"
)

type strategyRandom struct {
	candidatesSize int
	history        []string
	historySet     map[string]bool
}

func (s *strategyRandom) init() {
	s.history = nil
	s.historySet = map[string]bool{}
}

func (s *strategyRandom) GetTitleAndDescription() (string, string) {
	return "random", "Use random commands"
}

func (s *strategyRandom) GetCandidates() []string {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	var candidates []string
	candidateSet := map[string]bool{}
	for len(candidates) < s.candidatesSize && len(candidates)*2 < len(s.historySet) {
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

func (s *strategyRandom) AddHistoryRecord(record *records.EnrichedRecord) error {
	s.history = append([]string{record.CmdLine}, s.history...)
	s.historySet[record.CmdLine] = true
	return nil
}

func (s *strategyRandom) ResetHistory() error {
	s.init()
	return nil
}
