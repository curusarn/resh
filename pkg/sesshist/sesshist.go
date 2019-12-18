package sesshist

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/curusarn/resh/pkg/histfile"
	"github.com/curusarn/resh/pkg/histlist"
	"github.com/curusarn/resh/pkg/records"
)

// Dispatch Recall() calls to an apropriate session history (sesshist)
type Dispatch struct {
	sessions map[string]*sesshist
	mutex    sync.RWMutex

	history         *histfile.Histfile
	historyInitSize int
}

// NewDispatch creates a new sesshist.Dispatch and starts necessary gorutines
func NewDispatch(sessionsToInit chan records.Record, sessionsToDrop chan string,
	recordsToAdd chan records.Record, history *histfile.Histfile, historyInitSize int) *Dispatch {

	s := Dispatch{
		sessions:        map[string]*sesshist{},
		history:         history,
		historyInitSize: historyInitSize,
	}
	go s.sessionInitializer(sessionsToInit)
	go s.sessionDropper(sessionsToDrop)
	go s.recordAdder(recordsToAdd)
	return &s
}

func (s *Dispatch) sessionInitializer(sessionsToInit chan records.Record) {
	for {
		record := <-sessionsToInit
		log.Println("sesshist: got session to init - " + record.SessionID)
		s.initSession(record.SessionID, record.Shell)
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
		} else {
			// this inits session on RESH update
			s.checkSession(record.SessionID, record.Shell)
		}
		// TODO: we will need to handle part2 as well eventually
	}
}

func (s *Dispatch) checkSession(sessionID, shell string) {
	s.mutex.RLock()
	_, found := s.sessions[sessionID]
	s.mutex.RUnlock()
	if found == false {
		err := s.initSession(sessionID, shell)
		if err != nil {
			log.Println("sesshist: Error while checking session:", err)
		}
	}
}

// InitSession struct
func (s *Dispatch) initSession(sessionID, shell string) error {
	log.Println("sesshist: initializing session - " + sessionID)
	s.mutex.RLock()
	_, found := s.sessions[sessionID]
	s.mutex.RUnlock()

	if found == true {
		return errors.New("sesshist ERROR: Can't INIT already existing session " + sessionID)
	}

	log.Println("sesshist: loading history to populate session - " + sessionID)
	historyCmdLines := s.history.GetRecentCmdLines(shell, s.historyInitSize)

	s.mutex.Lock()
	defer s.mutex.Unlock()
	// init sesshist and populate it with history loaded from file
	s.sessions[sessionID] = &sesshist{
		recentCmdLines: historyCmdLines,
	}
	log.Println("sesshist: session init done - " + sessionID)
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
	log.Println("sesshist: Adding a record, RLocking main lock ...")
	s.mutex.RLock()
	log.Println("sesshist: Getting a session ...")
	session, found := s.sessions[sessionID]
	log.Println("sesshist: RUnlocking main lock ...")
	s.mutex.RUnlock()

	if found == false {
		log.Println("sesshist ERROR: addRecentRecord(): No session history for SessionID " + sessionID + " - creating session history.")
		s.initSession(sessionID, record.Shell)
		return s.addRecentRecord(sessionID, record)
	}
	log.Println("sesshist: RLocking session lock (w/ defer) ...")
	session.mutex.Lock()
	defer session.mutex.Unlock()
	session.recentRecords = append(session.recentRecords, record)
	session.recentCmdLines.AddCmdLine(record.CmdLine)
	log.Println("sesshist: record:", record.CmdLine, "; added to session:", sessionID,
		"; session len:", len(session.recentCmdLines.List), "; session len (records):", len(session.recentRecords))
	return nil
}

// Recall command from recent session history
func (s *Dispatch) Recall(sessionID string, histno int, prefix string) (string, error) {
	log.Println("sesshist - recall: RLocking main lock ...")
	s.mutex.RLock()
	log.Println("sesshist - recall: Getting session history struct ...")
	session, found := s.sessions[sessionID]
	s.mutex.RUnlock()

	if found == false {
		// TODO: propagate actual shell here so we can use it
		go s.initSession(sessionID, "bash")
		return "", errors.New("sesshist ERROR: No session history for SessionID " + sessionID + " - creating one ...")
	}
	log.Println("sesshist - recall: Locking session lock ...")
	session.mutex.Lock()
	defer session.mutex.Unlock()
	if prefix == "" {
		log.Println("sesshist - recall: Getting records by histno ...")
		return session.getRecordByHistno(histno)
	}
	log.Println("sesshist - recall: Searching for records by prefix ...")
	return session.searchRecordByPrefix(prefix, histno)
}

// Inspect commands in recent session history
func (s *Dispatch) Inspect(sessionID string, count int) ([]string, error) {
	prefix := ""
	log.Println("sesshist - inspect: RLocking main lock ...")
	s.mutex.RLock()
	log.Println("sesshist - inspect: Getting session history struct ...")
	session, found := s.sessions[sessionID]
	s.mutex.RUnlock()

	if found == false {
		// go s.initSession(sessionID)
		return nil, errors.New("sesshist ERROR: No session history for SessionID " + sessionID + " - should we create one?")
	}
	log.Println("sesshist - inspect: Locking session lock ...")
	session.mutex.Lock()
	defer session.mutex.Unlock()
	if prefix == "" {
		log.Println("sesshist - inspect: Getting records by histno ...")
		idx := len(session.recentCmdLines.List) - count
		if idx < 0 {
			idx = 0
		}
		return session.recentCmdLines.List[idx:], nil
	}
	log.Println("sesshist - inspect: Searching for records by prefix ... ERROR - Not implemented")
	return nil, errors.New("sesshist ERROR: Inspect - Searching for records by prefix Not implemented yet")
}

type sesshist struct {
	mutex          sync.Mutex
	recentRecords  []records.Record
	recentCmdLines histlist.Histlist
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
	index := len(s.recentCmdLines.List) - histno
	if index < 0 {
		return "", errors.New("sesshist ERROR: 'histno > number of commands in the session' (" + strconv.Itoa(len(s.recentCmdLines.List)) + ")")
	}
	return s.recentCmdLines.List[index], nil
}

func (s *sesshist) searchRecordByPrefix(prefix string, histno int) (string, error) {
	if histno == 0 {
		return "", errors.New("sesshist ERROR: 'histno == 0' is not a record from history")
	}
	if histno < 0 {
		return "", errors.New("sesshist ERROR: 'histno < 0' is a command from future (not supperted yet)")
	}
	index := len(s.recentCmdLines.List) - histno
	if index < 0 {
		return "", errors.New("sesshist ERROR: 'histno > number of commands in the session' (" + strconv.Itoa(len(s.recentCmdLines.List)) + ")")
	}
	cmdLines := []string{}
	for i := len(s.recentCmdLines.List) - 1; i >= 0; i-- {
		if strings.HasPrefix(s.recentCmdLines.List[i], prefix) {
			cmdLines = append(cmdLines, s.recentCmdLines.List[i])
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
