package histlist

import "go.uber.org/zap"

// Histlist is a deduplicated list of cmdLines
type Histlist struct {
	// TODO: I'm not excited about logger being passed here
	sugar *zap.SugaredLogger
	// list of commands lines (deduplicated)
	List []string
	// lookup: cmdLine -> last index
	LastIndex map[string]int
}

// New Histlist
func New(sugar *zap.SugaredLogger) Histlist {
	return Histlist{
		sugar:     sugar.With("component", "histlist"),
		LastIndex: make(map[string]int),
	}
}

// Copy Histlist
func Copy(hl Histlist) Histlist {
	newHl := New(hl.sugar)
	// copy list
	newHl.List = make([]string, len(hl.List))
	copy(newHl.List, hl.List)
	// copy map
	for k, v := range hl.LastIndex {
		newHl.LastIndex[k] = v
	}
	return newHl
}

// AddCmdLine to the histlist
func (h *Histlist) AddCmdLine(cmdLine string) {
	// lenBefore := len(h.list)
	// lookup
	idx, found := h.LastIndex[cmdLine]
	if found {
		// remove duplicate
		if cmdLine != h.List[idx] {
			h.sugar.DPanicw("Index key is different than actual cmd line in the list",
				"indexKeyCmdLine", cmdLine,
				"actualCmdLine", h.List[idx],
			)
		}
		h.List = append(h.List[:idx], h.List[idx+1:]...)
		// idx++
		for idx < len(h.List) {
			cmdLn := h.List[idx]
			h.LastIndex[cmdLn]--
			if idx != h.LastIndex[cmdLn] {
				h.sugar.DPanicw("Index position is different than actual position of the cmd line",
					"actualPosition", idx,
					"indexedPosition", h.LastIndex[cmdLn],
				)
			}
			idx++
		}
	}
	// update last index
	h.LastIndex[cmdLine] = len(h.List)
	// append new cmdline
	h.List = append(h.List, cmdLine)
	h.sugar.Debugw("Added cmdLine",
		"cmdLine", cmdLine,
		"historyLength", len(h.List),
	)
}

// AddHistlist contents of another histlist to this histlist
func (h *Histlist) AddHistlist(h2 Histlist) {
	for _, cmdLine := range h2.List {
		h.AddCmdLine(cmdLine)
	}
}
