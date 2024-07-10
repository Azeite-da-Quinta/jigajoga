package kvstore

import "github.com/Azeite-da-Quinta/jigajoga/libs/go/user"

// Payload retrieved from redis
type Payload struct {
	user.Token
	event string
}

// Type implements interface
func (p Payload) Type() string {
	return p.event
}

// User implements interface
func (p Payload) User() user.Token {
	return p.Token
}
