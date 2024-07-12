package bomb

// StateKind type
type StateKind uint8

// enum kinds of players
const (
	Lobby StateKind = iota
	Running
	Over
)
