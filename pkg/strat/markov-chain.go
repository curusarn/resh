package strat

import (
	"sort"
	"strconv"

	"github.com/curusarn/resh/pkg/records"
	"github.com/mb-14/gomarkov"
)

type MarkovChain struct {
	Order   int
	history []string
}

type strMarkEntry struct {
	cmdLine   string
	transProb float64
}

func (s *MarkovChain) Init() {
	s.history = nil
}

func (s *MarkovChain) GetTitleAndDescription() (string, string) {
	return "markov chain (order " + strconv.Itoa(s.Order) + ")", "Use markov chain to recommend commands"
}

func (s *MarkovChain) GetCandidates() []string {
	if len(s.history) < s.Order {
		return s.history
	}
	chain := gomarkov.NewChain(s.Order)

	chain.Add(s.history)

	cmdLinesSet := map[string]bool{}
	var entries []strMarkEntry
	for _, cmdLine := range s.history {
		if cmdLinesSet[cmdLine] {
			continue
		}
		cmdLinesSet[cmdLine] = true
		prob, _ := chain.TransitionProbability(cmdLine, s.history[len(s.history)-s.Order:])
		entries = append(entries, strMarkEntry{cmdLine: cmdLine, transProb: prob})
	}
	sort.Slice(entries, func(i int, j int) bool { return entries[i].transProb > entries[j].transProb })
	var hist []string
	for _, item := range entries {
		hist = append(hist, item.cmdLine)
	}
	// log.Println("################")
	// log.Println(s.history[len(s.history)-s.order:])
	// log.Println(" -> ")
	// x := math.Min(float64(len(hist)), 5)
	// log.Println(hist[:int(x)])
	// log.Println("################")
	return hist
}

func (s *MarkovChain) AddHistoryRecord(record *records.EnrichedRecord) error {
	s.history = append(s.history, record.CmdLine)
	// s.historySet[record.CmdLine] = true
	return nil
}

func (s *MarkovChain) ResetHistory() error {
	s.Init()
	return nil
}
