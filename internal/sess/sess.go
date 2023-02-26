package sess

// Session represents a session, used for sennding through channels when more than just ID is needed
type Session struct {
	ID  string
	PID int
}
