package user

import (
	"strconv"
	"sync/atomic"
)

// Mock implements the Token interface
type Mock struct {
	name   string
	roomID Identifier
	id     Identifier
}

// ID implements Token
func (m Mock) ID() Identifier {
	return m.id
}

// RoomID implements Token
func (m Mock) RoomID() Identifier {
	return m.roomID
}

// Name implements Token
func (m Mock) Name() string {
	return m.name
}

var counter atomic.Int64

// MockToken returns a dummy arbitrary Token
// with a user named bob and a "unique" id
// They all end up in the same room
func MockToken() Token {
	const (
		incr = 9
		base = 10
		room = 1234
	)
	id := counter.Add(incr)

	return Mock{
		roomID: room,
		id:     Identifier(id),
		name:   "bob" + strconv.FormatInt(id, base),
	}
}
