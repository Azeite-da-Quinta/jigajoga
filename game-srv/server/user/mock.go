package user

import (
	"strconv"
	"sync/atomic"
)

// Mock implements the Token interface
type Mock struct {
	roomID Identifier
	id     Identifier
	name   string
}

func (m Mock) ID() Identifier {
	return m.id
}

func (m Mock) RoomID() Identifier {
	return m.roomID
}

func (m Mock) Name() string {
	return m.name
}

var counter atomic.Int64

// MockToken returns a dummy arbitrary Token
// with a user named bob and a "unique" id
// They all end up in the same room
func MockToken() Token {
	id := strconv.FormatInt(counter.Add(9), 10)

	return Mock{
		roomID: "the-testing-room",
		id:     Identifier(id),
		name:   "bob" + id,
	}
}
