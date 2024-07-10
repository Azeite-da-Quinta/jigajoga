// Package comms for internal message passing
package comms

import (
	"strconv"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/libs/go/envelope"
	"github.com/Azeite-da-Quinta/jigajoga/libs/go/user"
)

// Message passed internally
type Message struct {
	Timestamp time.Time
	Content   string
	Target    int64
	Sender    int64
	Room      int64
}

// Fill extra fields
func (m *Message) Fill(u user.Token) {
	m.Sender = int64(u.ID())
	m.Room = int64(u.RoomID())
}

const base = 10

// Envelope converts to an envelope Message
func (m Message) Envelope() envelope.Message {
	return envelope.Message{
		From:    strconv.FormatInt(m.Sender, base),
		To:      strconv.FormatInt(m.Target, base),
		Content: m.Content,
	}
}

// FromEnvelope converts to internal Message
func FromEnvelope(e envelope.Message) (Message, error) {
	const (
		bitSize     = 64
		minFromSize = 0
	)

	var (
		sender int64
		err    error
	)
	if len(e.From) > minFromSize {
		sender, err = strconv.ParseInt(e.From, base, bitSize)
		if err != nil {
			return Message{}, err
		}
	}

	target, err := strconv.ParseInt(e.To, base, bitSize)
	if err != nil {
		return Message{}, err
	}

	return Message{
		Sender:  sender,
		Target:  target,
		Content: e.Content,
	}, nil
}
