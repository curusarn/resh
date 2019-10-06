package sesswatch

import (
	"log"
	"sync"
	"time"

	"github.com/curusarn/resh/pkg/records"
	"github.com/mitchellh/go-ps"
)

type sesswatch struct {
	sessionsToDrop []chan string
	sleepSeconds   uint

	watchedSessions map[string]bool
	mutex           sync.Mutex
}

// Go runs the session watcher - watches sessions and sends
func Go(sessionsToWatch chan records.Record, sessionsToDrop []chan string, sleepSeconds uint) {
	sw := sesswatch{sessionsToDrop: sessionsToDrop, sleepSeconds: sleepSeconds, watchedSessions: map[string]bool{}}
	go sw.waiter(sessionsToWatch)
}

func (s *sesswatch) waiter(sessionsToWatch chan records.Record) {
	for {
		func() {
			record := <-sessionsToWatch
			id := record.SessionID
			pid := record.SessionPID
			s.mutex.Lock()
			defer s.mutex.Unlock()
			if s.watchedSessions[id] == false {
				log.Println("sesswatch: start watching NEW session ~ pid:", id, "~", pid)
				s.watchedSessions[id] = true
				go s.watcher(id, pid)
			}
		}()
	}
}

func (s *sesswatch) watcher(sessionID string, sessionPID int) {
	for {
		time.Sleep(time.Duration(s.sleepSeconds) * time.Second)
		proc, err := ps.FindProcess(sessionPID)
		if err != nil {
			log.Println("sesswatch ERROR: error while finding process:", sessionPID)
		} else if proc == nil {
			log.Println("sesswatch: Dropping session ~ pid:", sessionID, "~", sessionPID)
			func() {
				s.mutex.Lock()
				defer s.mutex.Unlock()
				s.watchedSessions[sessionID] = false
			}()
			for _, ch := range s.sessionsToDrop {
				log.Println("sesswatch: sending 'drop session' message ...")
				ch <- sessionID
				log.Println("sesswatch: sending 'drop session' message DONE")
			}
			break
		}
	}
}
