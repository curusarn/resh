package syncconnector

import (
	"github.com/curusarn/resh/internal/recordint"
)

func (sc SyncConnector) write(collect chan recordint.Collect) {
	//for {
	//	func() {
	//		select {
	//		case rec := <-collect:
	//			part := "2"
	//			if rec.Rec.PartOne {
	//				part = "1"
	//			}
	//			sugar := h.sugar.With(
	//				"recordCmdLine", rec.Rec.CmdLine,
	//				"recordPart", part,
	//				"recordShell", rec.Shell,
	//			)
	//			sc.sugar.Debugw("Got record")
	//			h.sessionsMutex.Lock()
	//			defer h.sessionsMutex.Unlock()
	//
	//			// allows nested sessions to merge records properly
	//			mergeID := rec.SessionID + "_" + strconv.Itoa(rec.Shlvl)
	//			sugar = sc.sugar.With("mergeID", mergeID)
	//			if rec.Rec.PartOne {
	//				if _, found := h.sessions[mergeID]; found {
	//					msg := "Got another first part of the records before merging the previous one - overwriting!"
	//					if rec.Shell == "zsh" {
	//						sc.sugar.Warnw(msg)
	//					} else {
	//						sc.sugar.Infow(msg + " Unfortunately this is normal in bash, it can't be prevented.")
	//					}
	//				}
	//				h.sessions[mergeID] = rec
	//			} else {
	//				if part1, found := h.sessions[mergeID]; found == false {
	//					sc.sugar.Warnw("Got second part of record and nothing to merge it with - ignoring!")
	//				} else {
	//					delete(h.sessions, mergeID)
	//					go h.mergeAndWriteRecord(sugar, part1, rec)
	//				}
	//			}
	//		}
	//	}()
	//}
}
