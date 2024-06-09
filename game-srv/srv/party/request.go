package party

import (
	"context"
	"sync"

	"github.com/Azeite-da-Quinta/jigajoga/libs/user"
)

// Request incoming
type Request interface {
	Get() user.Token
}

// Join implements Request
type Join interface {
	Request
	Client
}

// Client in a room
type Client interface {
	user.Token
	Inbox() PostChan
	Cancel() context.CancelFunc
	Reply() chan JoinReply
}

// JoinReply contains the room's inbox and a WG to signal
// when a client is done writing in it
// TODO check if it would make sense to return an interface
type JoinReply struct {
	RoomInbox PostChan
	// TODO the Add(1) should be on the side of Router probably
	Wg *sync.WaitGroup
}

// Leave implements Request
type Leave interface {
	Request
	user.Token
}
