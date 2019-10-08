package sesshist

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/curusarn/resh/pkg/records"
)

// Dispatch Recall() calls to an apropriate session history (sesshist)
type Dispatch struct {
	sessions map[string]*sesshist
	mutex    sync.RWMutex
}

// NewDispatch creates a new sesshist.Dispatch and starts necessary gorutines
func NewDispatch(sessionsToInit chan records.Record, sessionsToDrop chan string, recordsToAdd chan records.Record) *Dispatch {
	s := Dispatch{sessions: map[string]*sesshist{}}
	go s.sessionInitializer(sessionsToInit)
	go s.sessionDropper(sessionsToDrop)
	go s.recordAdder(recordsToAdd)
	return &s
}

func (s *Dispatch) sessionInitializer(sessionsToInit chan records.Record) {
	for {
		record := <-sessionsToInit
		log.Println("sesshist: got session to init - " + record.SessionID)
		s.initSession(record.SessionID)
	}
}

func (s *Dispatch) sessionDropper(sessionsToDrop chan string) {
	for {
		sessionID := <-sessionsToDrop
		log.Println("sesshist: got session to drop - " + sessionID)
		s.dropSession(sessionID)
	}
}

func (s *Dispatch) recordAdder(recordsToAdd chan records.Record) {
	for {
		record := <-recordsToAdd
		if record.PartOne {
			log.Println("sesshist: got record to add - " + record.CmdLine)
			s.addRecentRecord(record.SessionID, record)
		}
		// TODO: we will need to handle part2 as well eventually
	}
}

// InitSession struct
func (s *Dispatch) initSession(sessionID string) error {
	s.mutex.RLock()
	_, found := s.sessions[sessionID]
	s.mutex.RUnlock()

	if found == true {
		return errors.New("sesshist ERROR: Can't INIT already existing session " + sessionID)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessions[sessionID] = &sesshist{}
	return nil
}

// DropSession struct
func (s *Dispatch) dropSession(sessionID string) error {
	s.mutex.RLock()
	_, found := s.sessions[sessionID]
	s.mutex.RUnlock()

	if found == false {
		return errors.New("sesshist ERROR: Can't DROP not existing session " + sessionID)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.sessions, sessionID)
	return nil
}

// AddRecent record to session
func (s *Dispatch) addRecentRecord(sessionID string, record records.Record) error {
	s.mutex.RLock()
	session, found := s.sessions[sessionID]
	s.mutex.RUnlock()

	if found == false {
		return errors.New("sesshist ERROR: No session history for SessionID " + sessionID + " (should we create one?)")
	}
	session.mutex.Lock()
	defer session.mutex.Unlock()
	session.recentRecords = append(session.recentRecords, record)
	// remove previous occurance of record
	for i := len(session.recentCmdLines) - 1; i >= 0; i-- {
		if session.recentCmdLines[i] == record.CmdLine {
			session.recentCmdLines = append(session.recentCmdLines[:i], session.recentCmdLines[i+1:]...)
		}
	}
	// append new record
	session.recentCmdLines = append(session.recentCmdLines, record.CmdLine)
	log.Println("sesshist: record:", record.CmdLine, "; added to session:", sessionID,
		"; session len:", len(session.recentCmdLines), "; session len w/ dups:", len(session.recentRecords))
	return nil
}

// Recall command from recent session history
func (s *Dispatch) Recall(sessionID string, histno int, prefix string) (string, error) {
	s.mutex.RLock()
	session, found := s.sessions[sessionID]
	s.mutex.RUnlock()

	if found == false {
		return "", errors.New("sesshist ERROR: No session history for SessionID " + sessionID)
	}
	if prefix == "" {
		session.mutex.Lock()
		defer session.mutex.Unlock()
		return session.getRecordByHistno(histno)
	}
	session.mutex.Lock()
	defer session.mutex.Unlock()
	return session.searchRecordByPrefix(prefix, histno)
}

type sesshist struct {
	recentRecords  []records.Record
	recentCmdLines []string // deduplicated
	// cmdLines map[string]int
	mutex sync.Mutex
}

func (s *sesshist) getRecordByHistno(histno int) (string, error) {
	// addRecords() appends records to the end of the slice
	// 	-> this func handles the indexing
	if histno == 0 {
		return "", errors.New("sesshist ERROR: 'histno == 0' is not a record from history")
	}
	if histno < 0 {
		return "", errors.New("sesshist ERROR: 'histno < 0' is a command from future (not supperted yet)")
	}
	index := len(s.recentCmdLines) - histno
	if index < 0 {
		return "", errors.New("sesshist ERROR: 'histno > number of commands in the session' (" + strconv.Itoa(len(s.recentCmdLines)) + ")")
	}
	return s.recentCmdLines[index], nil
}

func (s *sesshist) searchRecordByPrefix(prefix string, histno int) (string, error) {
	if histno == 0 {
		return "", errors.New("sesshist ERROR: 'histno == 0' is not a record from history")
	}
	if histno < 0 {
		return "", errors.New("sesshist ERROR: 'histno < 0' is a command from future (not supperted yet)")
	}
	index := len(s.recentCmdLines) - histno
	if index < 0 {
		return "", errors.New("sesshist ERROR: 'histno > number of commands in the session' (" + strconv.Itoa(len(s.recentCmdLines)) + ")")
	}
	cmdLines := []string{}
	for i := len(s.recentCmdLines) - 1; i >= 0; i-- {
		if strings.HasPrefix(s.recentCmdLines[i], prefix) {
			cmdLines = append(cmdLines, s.recentCmdLines[i])
			if len(cmdLines) >= histno {
				break
			}
		}
	}
	if len(cmdLines) < histno {
		return "", errors.New("sesshist ERROR: 'histno > number of commands matching with given prefix' (" + strconv.Itoa(len(cmdLines)) + ")")
	}
	return cmdLines[histno-1], nil
}
