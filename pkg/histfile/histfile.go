package histfile

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"strconv"
	"sync"

	"github.com/curusarn/resh/pkg/histcli"
	"github.com/curusarn/resh/pkg/histlist"
	"github.com/curusarn/resh/pkg/records"
)

// Histfile writes records to histfile
type Histfile struct {
	sessionsMutex sync.Mutex
	sessions      map[string]records.Record
	historyPath   string

	recentMutex   sync.Mutex
	recentRecords []records.Record

	// NOTE: we have separate histories which only differ if there was not enough resh_history
	//			resh_history itself is common for both bash and zsh
	bashCmdLines histlist.Histlist
	zshCmdLines  histlist.Histlist

	fullRecords histcli.Histcli
}

// New creates new histfile and runs its gorutines
func New(input chan records.Record, sessionsToDrop chan string,
	reshHistoryPath string, bashHistoryPath string, zshHistoryPath string,
	maxInitHistSize int, minInitHistSizeKB int,
	signals chan os.Signal, shutdownDone chan string) *Histfile {

	hf := Histfile{
		sessions:     map[string]records.Record{},
		historyPath:  reshHistoryPath,
		bashCmdLines: histlist.New(),
		zshCmdLines:  histlist.New(),
		fullRecords:  histcli.New(),
	}
	go hf.loadHistory(bashHistoryPath, zshHistoryPath, maxInitHistSize, minInitHistSizeKB)
	go hf.writer(input, signals, shutdownDone)
	go hf.sessionGC(sessionsToDrop)
	go hf.loadFullRecords()
	return &hf
}

// load records from resh history, reverse, enrich and save
func (h *Histfile) loadFullRecords() {
	recs := records.LoadFromFile(h.historyPath, math.MaxInt32)
	for i := len(recs) - 1; i >= 0; i-- {
		rec := recs[i]
		h.fullRecords.AddRecord(rec)
	}
}

// loadsHistory from resh_history and if there is not enough of it also load native shell histories
func (h *Histfile) loadHistory(bashHistoryPath, zshHistoryPath string, maxInitHistSize, minInitHistSizeKB int) {
	h.recentMutex.Lock()
	defer h.recentMutex.Unlock()
	log.Println("histfile: Checking if resh_history is large enough ...")
	fi, err := os.Stat(h.historyPath)
	var size int
	if err != nil {
		log.Println("histfile ERROR: failed to stat resh_history file:", err)
	} else {
		size = int(fi.Size())
	}
	useNativeHistories := false
	if size/1024 < minInitHistSizeKB {
		useNativeHistories = true
		log.Println("histfile WARN: resh_history is too small - loading native bash and zsh history ...")
		h.bashCmdLines = records.LoadCmdLinesFromBashFile(bashHistoryPath)
		log.Println("histfile: bash history loaded - cmdLine count:", len(h.bashCmdLines.List))
		h.zshCmdLines = records.LoadCmdLinesFromZshFile(zshHistoryPath)
		log.Println("histfile: zsh history loaded - cmdLine count:", len(h.zshCmdLines.List))
		// no maxInitHistSize when using native histories
		maxInitHistSize = math.MaxInt32
	}
	log.Println("histfile: Loading resh history from file ...")
	reshCmdLines := histlist.New()
	// NOTE: keeping this weird interface for now because we might use it in the future
	//			when we only load bash or zsh history
	records.LoadCmdLinesFromFile(&reshCmdLines, h.historyPath, maxInitHistSize)
	log.Println("histfile: resh history loaded - cmdLine count:", len(reshCmdLines.List))
	if useNativeHistories == false {
		h.bashCmdLines = reshCmdLines
		h.zshCmdLines = histlist.Copy(reshCmdLines)
		return
	}
	h.bashCmdLines.AddHistlist(reshCmdLines)
	log.Println("histfile: bash history + resh history - cmdLine count:", len(h.bashCmdLines.List))
	h.zshCmdLines.AddHistlist(reshCmdLines)
	log.Println("histfile: zsh history + resh history - cmdLine count:", len(h.zshCmdLines.List))
}

// sessionGC reads sessionIDs from channel and deletes them from histfile struct
func (h *Histfile) sessionGC(sessionsToDrop chan string) {
	for {
		func() {
			session := <-sessionsToDrop
			log.Println("histfile: got session to drop", session)
			h.sessionsMutex.Lock()
			defer h.sessionsMutex.Unlock()
			if part1, found := h.sessions[session]; found == true {
				log.Println("histfile: Dropping session:", session)
				delete(h.sessions, session)
				go writeRecord(part1, h.historyPath)
			} else {
				log.Println("histfile: No hanging parts for session:", session)
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
				h.sessionsMutex.Lock()
				defer h.sessionsMutex.Unlock()

				// allows nested sessions to merge records properly
				mergeID := record.SessionID + "_" + strconv.Itoa(record.Shlvl)
				if record.PartOne {
					if _, found := h.sessions[mergeID]; found {
						log.Println("histfile WARN: Got another first part of the records before merging the previous one - overwriting! " +
							"(this happens in bash because bash-preexec runs when it's not supposed to)")
					}
					h.sessions[mergeID] = record
				} else {
					if part1, found := h.sessions[mergeID]; found == false {
						log.Println("histfile ERROR: Got second part of records and nothing to merge it with - ignoring! (mergeID:", mergeID, ")")
					} else {
						delete(h.sessions, mergeID)
						go h.mergeAndWriteRecord(part1, record)
					}
				}
			case sig := <-signals:
				log.Println("histfile: Got signal " + sig.String())
				h.sessionsMutex.Lock()
				defer h.sessionsMutex.Unlock()
				log.Println("histfile DEBUG: Unlocked mutex")

				for sessID, record := range h.sessions {
					log.Panicln("histfile WARN: Writing incomplete record for session " + sessID)
					h.writeRecord(record)
				}
				log.Println("histfile DEBUG: Shutdown success")
				shutdownDone <- "histfile"
				return
			}
		}()
	}
}

func (h *Histfile) writeRecord(part1 records.Record) {
	writeRecord(part1, h.historyPath)
}

func (h *Histfile) mergeAndWriteRecord(part1, part2 records.Record) {
	err := part1.Merge(part2)
	if err != nil {
		log.Println("Error while merging", err)
		return
	}

	func() {
		h.recentMutex.Lock()
		defer h.recentMutex.Unlock()
		h.recentRecords = append(h.recentRecords, part1)
		cmdLine := part1.CmdLine
		h.bashCmdLines.AddCmdLine(cmdLine)
		h.zshCmdLines.AddCmdLine(cmdLine)
		h.fullRecords.AddRecord(part1)
	}()

	writeRecord(part1, h.historyPath)
}

func writeRecord(rec records.Record, outputPath string) {
	recJSON, err := json.Marshal(rec)
	if err != nil {
		log.Println("Marshalling error", err)
		return
	}
	f, err := os.OpenFile(outputPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Could not open file", err)
		return
	}
	defer f.Close()
	_, err = f.Write(append(recJSON, []byte("\n")...))
	if err != nil {
		log.Printf("Error while writing: %v, %s\n", rec, err)
		return
	}
}

// GetRecentCmdLines returns recent cmdLines
func (h *Histfile) GetRecentCmdLines(shell string, limit int) histlist.Histlist {
	// NOTE: limit does nothing atm
	h.recentMutex.Lock()
	defer h.recentMutex.Unlock()
	log.Println("histfile: History requested ...")
	var hl histlist.Histlist
	if shell == "bash" {
		hl = histlist.Copy(h.bashCmdLines)
		log.Println("histfile: history copied (bash) - cmdLine count:", len(hl.List))
		return hl
	}
	if shell != "zsh" {
		log.Println("histfile ERROR: Unknown shell: ", shell)
	}
	hl = histlist.Copy(h.zshCmdLines)
	log.Println("histfile: history copied (zsh) - cmdLine count:", len(hl.List))
	return hl
}

// DumpRecords returns enriched records
func (h *Histfile) DumpRecords() histcli.Histcli {
	// don't forget locks in the future
	return h.fullRecords
}
