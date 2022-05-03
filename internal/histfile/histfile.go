package histfile

import (
	"encoding/json"
	"math"
	"os"
	"strconv"
	"sync"

	"github.com/curusarn/resh/internal/histcli"
	"github.com/curusarn/resh/internal/histlist"
	"github.com/curusarn/resh/internal/records"
	"go.uber.org/zap"
)

// Histfile writes records to histfile
type Histfile struct {
	sugar *zap.SugaredLogger

	sessionsMutex sync.Mutex
	sessions      map[string]records.Record
	historyPath   string

	recentMutex   sync.Mutex
	recentRecords []records.Record

	// NOTE: we have separate histories which only differ if there was not enough resh_history
	//			resh_history itself is common for both bash and zsh
	bashCmdLines histlist.Histlist
	zshCmdLines  histlist.Histlist

	cliRecords histcli.Histcli
}

// New creates new histfile and runs its gorutines
func New(sugar *zap.SugaredLogger, input chan records.Record, sessionsToDrop chan string,
	reshHistoryPath string, bashHistoryPath string, zshHistoryPath string,
	maxInitHistSize int, minInitHistSizeKB int,
	signals chan os.Signal, shutdownDone chan string) *Histfile {

	hf := Histfile{
		sugar:        sugar.With("module", "histfile"),
		sessions:     map[string]records.Record{},
		historyPath:  reshHistoryPath,
		bashCmdLines: histlist.New(sugar),
		zshCmdLines:  histlist.New(sugar),
		cliRecords:   histcli.New(),
	}
	go hf.loadHistory(bashHistoryPath, zshHistoryPath, maxInitHistSize, minInitHistSizeKB)
	go hf.writer(input, signals, shutdownDone)
	go hf.sessionGC(sessionsToDrop)
	return &hf
}

// load records from resh history, reverse, enrich and save
func (h *Histfile) loadCliRecords(recs []records.Record) {
	for _, cmdline := range h.bashCmdLines.List {
		h.cliRecords.AddCmdLine(cmdline)
	}
	for _, cmdline := range h.zshCmdLines.List {
		h.cliRecords.AddCmdLine(cmdline)
	}
	for i := len(recs) - 1; i >= 0; i-- {
		rec := recs[i]
		h.cliRecords.AddRecord(rec)
	}
	h.sugar.Infow("Resh history loaded",
		"historyRecordsCount", len(h.cliRecords.List),
	)
}

// loadsHistory from resh_history and if there is not enough of it also load native shell histories
func (h *Histfile) loadHistory(bashHistoryPath, zshHistoryPath string, maxInitHistSize, minInitHistSizeKB int) {
	h.recentMutex.Lock()
	defer h.recentMutex.Unlock()
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
	history := records.LoadFromFile(h.sugar, h.historyPath)
	h.sugar.Infow("Resh history loaded from file",
		"historyFile", h.historyPath,
		"recordCount", len(history),
	)
	go h.loadCliRecords(history)
	// NOTE: keeping this weird interface for now because we might use it in the future
	//			when we only load bash or zsh history
	reshCmdLines := loadCmdLines(h.sugar, history)
	h.sugar.Infow("Resh history loaded and processed",
		"recordCount", len(reshCmdLines.List),
	)
	if useNativeHistories == false {
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
				go writeRecord(sugar, part1, h.historyPath)
			} else {
				sugar.Infow("No hanging parts for session - nothing to drop")
			}
		}()
	}
}

// writer reads records from channel, merges them and writes them to file
func (h *Histfile) writer(input chan records.Record, signals chan os.Signal, shutdownDone chan string) {
	for {
		func() {
			select {
			case record := <-input:
				part := "2"
				if record.PartOne {
					part = "1"
				}
				sugar := h.sugar.With(
					"recordCmdLine", record.CmdLine,
					"recordPart", part,
					"recordShell", record.Shell,
				)
				sugar.Debugw("Got record")
				h.sessionsMutex.Lock()
				defer h.sessionsMutex.Unlock()

				// allows nested sessions to merge records properly
				mergeID := record.SessionID + "_" + strconv.Itoa(record.Shlvl)
				sugar = sugar.With("mergeID", mergeID)
				if record.PartOne {
					if _, found := h.sessions[mergeID]; found {
						msg := "Got another first part of the records before merging the previous one - overwriting!"
						if record.Shell == "zsh" {
							sugar.Warnw(msg)
						} else {
							sugar.Infow(msg + " Unfortunately this is normal in bash, it can't be prevented.")
						}
					}
					h.sessions[mergeID] = record
				} else {
					if part1, found := h.sessions[mergeID]; found == false {
						sugar.Warnw("Got second part of record and nothing to merge it with - ignoring!")
					} else {
						delete(h.sessions, mergeID)
						go h.mergeAndWriteRecord(sugar, part1, record)
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

				for sessID, record := range h.sessions {
					sugar.Warnw("Writing incomplete record for session",
						"sessionID", sessID,
					)
					h.writeRecord(sugar, record)
				}
				sugar.Debugw("Shutdown successful")
				shutdownDone <- "histfile"
				return
			}
		}()
	}
}

func (h *Histfile) writeRecord(sugar *zap.SugaredLogger, part1 records.Record) {
	writeRecord(sugar, part1, h.historyPath)
}

func (h *Histfile) mergeAndWriteRecord(sugar *zap.SugaredLogger, part1, part2 records.Record) {
	err := part1.Merge(part2)
	if err != nil {
		sugar.Errorw("Error while merging records", "error", err)
		return
	}

	func() {
		h.recentMutex.Lock()
		defer h.recentMutex.Unlock()
		h.recentRecords = append(h.recentRecords, part1)
		cmdLine := part1.CmdLine
		h.bashCmdLines.AddCmdLine(cmdLine)
		h.zshCmdLines.AddCmdLine(cmdLine)
		h.cliRecords.AddRecord(part1)
	}()

	writeRecord(sugar, part1, h.historyPath)
}

func writeRecord(sugar *zap.SugaredLogger, rec records.Record, outputPath string) {
	recJSON, err := json.Marshal(rec)
	if err != nil {
		sugar.Errorw("Marshalling error", "error", err)
		return
	}
	f, err := os.OpenFile(outputPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		sugar.Errorw("Could not open file", "error", err)
		return
	}
	defer f.Close()
	_, err = f.Write(append(recJSON, []byte("\n")...))
	if err != nil {
		sugar.Errorw("Error while writing record",
			"recordRaw", rec,
			"error", err,
		)
		return
	}
}

// DumpCliRecords returns enriched records
func (h *Histfile) DumpCliRecords() histcli.Histcli {
	// don't forget locks in the future
	return h.cliRecords
}

func loadCmdLines(sugar *zap.SugaredLogger, recs []records.Record) histlist.Histlist {
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
