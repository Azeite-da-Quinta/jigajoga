// Package event contains incoming notifications for parties
//
//revive:disable:max-public-structs
package event

import (
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/party/comms"
	"github.com/Azeite-da-Quinta/jigajoga/libs/user"
)

// InboxChan is a read only message channel
type InboxChan <-chan comms.Message

// PostChan is a write only message channel
type PostChan chan<- comms.Message

// Query incoming
type Query interface {
	Envelope
}

// Envelope that you format to be transported
type Envelope interface {
	Topic() string
	Serialize() ([]byte, error)
}

// Join implements Event
type Join interface {
	Query
	Client
}

// Leave implements Event
type Leave interface {
	Query
	user.Token
}

// Client in a room
type Client interface {
	user.Token
	Reply() chan Reply
}

// Reply is the expected reply from Room
type Reply interface {
	// TODO Probably useless
	OriginID() int64
	// write only - send to room
	Room() PostChan
	// ready only - receive from room
	Inbox() InboxChan
}
