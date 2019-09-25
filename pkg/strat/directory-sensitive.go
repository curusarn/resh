package strat

import "github.com/curusarn/resh/pkg/records"

type DirectorySensitive struct {
	history map[string][]string
	lastPwd string
}

func (s *DirectorySensitive) Init() {
	s.history = map[string][]string{}
}

func (s *DirectorySensitive) GetTitleAndDescription() (string, string) {
	return "directory sensitive (recent)", "Use recent commands executed is the same directory"
}

func (s *DirectorySensitive) GetCandidates() []string {
	return s.history[s.lastPwd]
}

func (s *DirectorySensitive) AddHistoryRecord(record *records.EnrichedRecord) error {
	// work on history for PWD
	pwd := record.Pwd
	// remove previous occurance of record
	for i, cmd := range s.history[pwd] {
		if cmd == record.CmdLine {
			s.history[pwd] = append(s.history[pwd][:i], s.history[pwd][i+1:]...)
		}
	}
	// append new record
	s.history[pwd] = append([]string{record.CmdLine}, s.history[pwd]...)
	s.lastPwd = record.PwdAfter
	return nil
}

func (s *DirectorySensitive) ResetHistory() error {
	s.Init()
	s.history = map[string][]string{}
	return nil
}
