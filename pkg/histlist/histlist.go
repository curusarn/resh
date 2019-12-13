package histlist

// Histlist is a deduplicated list of cmdLines
type Histlist struct {
	// list of commands lines (deduplicated)
	List []string
	// lookup: cmdLine -> last index
	LastIndex map[string]int
}

// New Histlist
func New() Histlist {
	return Histlist{LastIndex: make(map[string]int)}
}

// Copy Histlist
func Copy(hl Histlist) Histlist {
	newHl := New()
	// copy list
	newHl.List = make([]string, len(hl.List))
	copy(newHl.List, hl.List)
	// copy map
	for k, v := range hl.LastIndex {
		newHl.LastIndex[k] = v
	}
	return newHl
}
