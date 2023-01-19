package histfile

import (
	"math"
	"os"
	"strconv"
	"sync"

	"github.com/curusarn/resh/internal/histcli"
	"github.com/curusarn/resh/internal/histlist"
	"github.com/curusarn/resh/internal/recio"
	"github.com/curusarn/resh/internal/recordint"
	"github.com/curusarn/resh/internal/records"
	"github.com/curusarn/resh/internal/recutil"
	"github.com/curusarn/resh/record"
	"go.uber.org/zap"
)

// TODO: get rid of histfile - use histio instead
// Histfile writes records to histfile
type Histfile struct {
	sugar *zap.SugaredLogger

	sessionsMutex sync.Mutex
	sessions      map[string]recordint.Collect
	historyPath   string

	// NOTE: we have separate histories which only differ if there was not enough resh_history
	//			resh_history itself is common for both bash and zsh
	bashCmdLines histlist.Histlist
	zshCmdLines  histlist.Histlist

	cliRecords histcli.Histcli

	rio *recio.RecIO
}

// New creates new histfile and runs its goroutines
func New(sugar *zap.SugaredLogger, input chan recordint.Collect, sessionsToDrop chan string,
	reshHistoryPath string, bashHistoryPath string, zshHistoryPath string,
	maxInitHistSize int, minInitHistSizeKB int,
	signals chan os.Signal, shutdownDone chan string) *Histfile {

	rio := recio.New(sugar.With("module", "histfile"))
	hf := Histfile{
		sugar:        sugar.With("module", "histfile"),
		sessions:     map[string]recordint.Collect{},
		historyPath:  reshHistoryPath,
		bashCmdLines: histlist.New(sugar),
		zshCmdLines:  histlist.New(sugar),
		cliRecords:   histcli.New(sugar),
		rio:          &rio,
	}
	go hf.loadHistory(bashHistoryPath, zshHistoryPath, maxInitHistSize, minInitHistSizeKB)
	go hf.writer(input, signals, shutdownDone)
	go hf.sessionGC(sessionsToDrop)
	return &hf
}

// load records from resh history, reverse, enrich and save
func (h *Histfile) loadCliRecords(recs []record.V1) {
	for _, cmdline := range h.bashCmdLines.List {
		h.cliRecords.AddCmdLine(cmdline)
	}
	for _, cmdline := range h.zshCmdLines.List {
		h.cliRecords.AddCmdLine(cmdline)
	}
	for i := len(recs) - 1; i >= 0; i-- {
		rec := recs[i]
		h.cliRecords.AddRecord(&rec)
	}
	h.sugar.Infow("Resh history loaded",
		"historyRecordsCount", len(h.cliRecords.List),
	)
}

// loadsHistory from resh_history and if there is not enough of it also load native shell histories
func (h *Histfile) loadHistory(bashHistoryPath, zshHistoryPath string, maxInitHistSize, minInitHistSizeKB int) {
	h.sugar.Infow("Checking if resh_history is large enough ...")
	fi, err := os.Stat(h.historyPath)
	var size int
	if err != nil {
		h.sugar.Errorw("Failed to stat resh_history file", "error", err)
	} else {
		size = int(fi.Size())
	}
	useNativeHistories := false
	if size/1024 < minInitHistSizeKB {
		useNativeHistories = true
		h.sugar.Warnw("Resh_history is too small - loading native bash and zsh history ...")
		h.bashCmdLines = records.LoadCmdLinesFromBashFile(h.sugar, bashHistoryPath)
		h.sugar.Infow("Bash history loaded", "cmdLineCount", len(h.bashCmdLines.List))
		h.zshCmdLines = records.LoadCmdLinesFromZshFile(h.sugar, zshHistoryPath)
		h.sugar.Infow("Zsh history loaded", "cmdLineCount", len(h.zshCmdLines.List))
		// no maxInitHistSize when using native histories
		maxInitHistSize = math.MaxInt32
	}
	h.sugar.Debugw("Loading resh history from file ...",
		"historyFile", h.historyPath,
	)
	history, err := h.rio.ReadAndFixFile(h.historyPath, 3)
	if err != nil {
		h.sugar.Fatalf("Failed to read history file: %v", err)
	}
	h.sugar.Infow("Resh history loaded from file",
		"historyFile", h.historyPath,
		"recordCount", len(history),
	)
	go h.loadCliRecords(history)
	// NOTE: keeping this weird interface for now because we might use it in the future
	//       when we only load bash or zsh history
	reshCmdLines := loadCmdLines(h.sugar, history)
	h.sugar.Infow("Resh history loaded and processed",
		"recordCount", len(reshCmdLines.List),
	)
	if !useNativeHistories {
		h.bashCmdLines = reshCmdLines
		h.zshCmdLines = histlist.Copy(reshCmdLines)
		return
	}
	h.bashCmdLines.AddHistlist(reshCmdLines)
	h.sugar.Infow("Processed bash history and resh history together", "cmdLinecount", len(h.bashCmdLines.List))
	h.zshCmdLines.AddHistlist(reshCmdLines)
	h.sugar.Infow("Processed zsh history and resh history together", "cmdLineCount", len(h.zshCmdLines.List))
}

// sessionGC reads sessionIDs from channel and deletes them from histfile struct
func (h *Histfile) sessionGC(sessionsToDrop chan string) {
	for {
		func() {
			session := <-sessionsToDrop
			sugar := h.sugar.With("sessionID", session)
			sugar.Debugw("Got session to drop")
			h.sessionsMutex.Lock()
			defer h.sessionsMutex.Unlock()
			if part1, found := h.sessions[session]; found == true {
				sugar.Infow("Dropping session")
				delete(h.sessions, session)
				go h.rio.AppendToFile(h.historyPath, []record.V1{part1.Rec})
			} else {
				sugar.Infow("No hanging parts for session - nothing to drop")
			}
		}()
	}
}

// writer reads records from channel, merges them and writes them to file
func (h *Histfile) writer(collect chan recordint.Collect, signals chan os.Signal, shutdownDone chan string) {
	for {
		func() {
			select {
			case rec := <-collect:
				part := "2"
				if rec.Rec.PartOne {
					part = "1"
				}
				sugar := h.sugar.With(
					"recordCmdLine", rec.Rec.CmdLine,
					"recordPart", part,
					"recordShell", rec.Shell,
				)
				sugar.Debugw("Got record")
				h.sessionsMutex.Lock()
				defer h.sessionsMutex.Unlock()

				// allows nested sessions to merge records properly
				mergeID := rec.SessionID + "_" + strconv.Itoa(rec.Shlvl)
				sugar = sugar.With("mergeID", mergeID)
				if rec.Rec.PartOne {
					if _, found := h.sessions[mergeID]; found {
						msg := "Got another first part of the records before merging the previous one - overwriting!"
						if rec.Shell == "zsh" {
							sugar.Warnw(msg)
						} else {
							sugar.Infow(msg + " Unfortunately this is normal in bash, it can't be prevented.")
						}
					}
					h.sessions[mergeID] = rec
				} else {
					if part1, found := h.sessions[mergeID]; found == false {
						sugar.Warnw("Got second part of record and nothing to merge it with - ignoring!")
					} else {
						delete(h.sessions, mergeID)
						go h.mergeAndWriteRecord(sugar, part1, rec)
					}
				}
			case sig := <-signals:
				sugar := h.sugar.With(
					"signal", sig.String(),
				)
				sugar.Infow("Got signal")
				h.sessionsMutex.Lock()
				defer h.sessionsMutex.Unlock()
				sugar.Debugw("Unlocked mutex")

				for sessID, rec := range h.sessions {
					sugar.Warnw("Writing incomplete record for session",
						"sessionID", sessID,
					)
					h.writeRecord(sugar, rec.Rec)
				}
				sugar.Debugw("Shutdown successful")
				shutdownDone <- "histfile"
				return
			}
		}()
	}
}

func (h *Histfile) writeRecord(sugar *zap.SugaredLogger, rec record.V1) {
	h.rio.AppendToFile(h.historyPath, []record.V1{rec})
}

func (h *Histfile) mergeAndWriteRecord(sugar *zap.SugaredLogger, part1 recordint.Collect, part2 recordint.Collect) {
	rec, err := recutil.Merge(&part1, &part2)
	if err != nil {
		sugar.Errorw("Error while merging records", "error", err)
		return
	}

	recV1 := record.V1(rec)
	func() {
		cmdLine := rec.CmdLine
		h.bashCmdLines.AddCmdLine(cmdLine)
		h.zshCmdLines.AddCmdLine(cmdLine)
		h.cliRecords.AddRecord(&recV1)
	}()

	h.rio.AppendToFile(h.historyPath, []record.V1{recV1})
}

// TODO: use errors in RecIO
// func writeRecord(sugar *zap.SugaredLogger, rec record.V1, outputPath string) {
// 	recJSON, err := json.Marshal(rec)
// 	if err != nil {
// 		sugar.Errorw("Marshalling error", "error", err)
// 		return
// 	}
// 	f, err := os.OpenFile(outputPath,
// 		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		sugar.Errorw("Could not open file", "error", err)
// 		return
// 	}
// 	defer f.Close()
// 	_, err = f.Write(append(recJSON, []byte("\n")...))
// 	if err != nil {
// 		sugar.Errorw("Error while writing record",
// 			"recordRaw", rec,
// 			"error", err,
// 		)
// 		return
// 	}
// }

// DumpCliRecords returns enriched records
func (h *Histfile) DumpCliRecords() histcli.Histcli {
	// don't forget locks in the future
	return h.cliRecords
}

func loadCmdLines(sugar *zap.SugaredLogger, recs []record.V1) histlist.Histlist {
	hl := histlist.New(sugar)
	// go from bottom and deduplicate
	var cmdLines []string
	cmdLinesSet := map[string]bool{}
	for i := len(recs) - 1; i >= 0; i-- {
		cmdLine := recs[i].CmdLine
		if cmdLinesSet[cmdLine] {
			continue
		}
		cmdLinesSet[cmdLine] = true
		cmdLines = append([]string{cmdLine}, cmdLines...)
		// if len(cmdLines) > limit {
		// 	break
		// }
	}
	// add everything to histlist
	for _, cmdLine := range cmdLines {
		hl.AddCmdLine(cmdLine)
	}
	return hl
}
