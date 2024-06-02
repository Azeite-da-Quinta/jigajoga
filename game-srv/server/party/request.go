package party

import "github.com/Azeite-da-Quinta/jigajoga/game-srv/server/user"

// Request incoming
type Request interface {
	Get() user.Token
}

// Join implements Request
// contains the client's inbox channel
// and a reply channel to provide the room's inbox
type Join struct {
	user.Token
	ClientInbox PostChan
	// Reply in this channel with the room's write channel
	ReplyRoom chan PostChan
}

func (j Join) Get() user.Token {
	return j.Token
}

// Leave implements Request
type Leave struct {
	user.Token
}

func (l Leave) Get() user.Token {
	return l.Token
}
