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

// JoinReply
type JoinReply struct {
	RoomInbox PostChan
	Wg        *sync.WaitGroup
}

// Leave implements Request
type Leave interface {
	Request
	user.Token
}
