package sesswatch

import (
	"sync"
	"time"

	"github.com/curusarn/resh/internal/recordint"
	"github.com/mitchellh/go-ps"
	"go.uber.org/zap"
)

type sesswatch struct {
	sugar *zap.SugaredLogger

	sessionsToDrop []chan string
	sleepSeconds   uint

	watchedSessions map[string]bool
	mutex           sync.Mutex
}

// Go runs the session watcher - watches sessions and sends
func Go(sugar *zap.SugaredLogger,
	sessionsToWatch chan recordint.SessionInit, sessionsToWatchRecords chan recordint.Collect,
	sessionsToDrop []chan string, sleepSeconds uint) {

	sw := sesswatch{
		sugar:           sugar.With("module", "sesswatch"),
		sessionsToDrop:  sessionsToDrop,
		sleepSeconds:    sleepSeconds,
		watchedSessions: map[string]bool{},
	}
	go sw.waiter(sessionsToWatch, sessionsToWatchRecords)
}

func (s *sesswatch) waiter(sessionsToWatch chan recordint.SessionInit, sessionsToWatchRecords chan recordint.Collect) {
	for {
		func() {
			select {
			case rec := <-sessionsToWatch:
				// normal way to start watching a session
				id := rec.SessionID
				pid := rec.SessionPID
				sugar := s.sugar.With(
					"sessionID", rec.SessionID,
					"sessionPID", rec.SessionPID,
				)
				s.mutex.Lock()
				defer s.mutex.Unlock()
				if s.watchedSessions[id] == false {
					sugar.Infow("Starting watching new session")
					s.watchedSessions[id] = true
					go s.watcher(sugar, id, pid)
				}
			case rec := <-sessionsToWatchRecords:
				// additional safety - watch sessions that were never properly initialized
				id := rec.SessionID
				pid := rec.SessionPID
				sugar := s.sugar.With(
					"sessionID", rec.SessionID,
					"sessionPID", rec.SessionPID,
				)
				s.mutex.Lock()
				defer s.mutex.Unlock()
				if s.watchedSessions[id] == false {
					sugar.Warnw("Starting watching new session based on '/record'")
					s.watchedSessions[id] = true
					go s.watcher(sugar, id, pid)
				}
			}
		}()
	}
}

func (s *sesswatch) watcher(sugar *zap.SugaredLogger, sessionID string, sessionPID int) {
	for {
		time.Sleep(time.Duration(s.sleepSeconds) * time.Second)
		proc, err := ps.FindProcess(sessionPID)
		if err != nil {
			sugar.Errorw("Error while finding process", "error", err)
		} else if proc == nil {
			sugar.Infow("Dropping session")
			func() {
				s.mutex.Lock()
				defer s.mutex.Unlock()
				s.watchedSessions[sessionID] = false
			}()
			for _, ch := range s.sessionsToDrop {
				sugar.Debugw("Sending 'drop session' message ...")
				ch <- sessionID
				sugar.Debugw("Sending 'drop session' message DONE")
			}
			break
		}
	}
}
