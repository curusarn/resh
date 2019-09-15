package main

import "github.com/curusarn/resh/common"

type strategyDirectorySensitive struct {
	history map[string][]string
	lastPwd string
}

func (s *strategyDirectorySensitive) init() {
	s.history = map[string][]string{}
}

func (s *strategyDirectorySensitive) GetTitleAndDescription() (string, string) {
	return "directory sensitive (recent)", "Use recent commands executed is the same directory"
}

func (s *strategyDirectorySensitive) GetCandidates() []string {
	return s.history[s.lastPwd]
}

func (s *strategyDirectorySensitive) AddHistoryRecord(record *common.EnrichedRecord) error {
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

func (s *strategyDirectorySensitive) ResetHistory() error {
	s.history = map[string][]string{}
	return nil
}
