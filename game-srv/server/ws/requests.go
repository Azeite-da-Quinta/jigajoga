package ws

import (
	"context"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/server/party"
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/server/user"
)

type join struct {
	user.Token
	inbox party.PostChan
	// Reply in this channel with the room's write channel
	reply chan party.JoinReply
	// client's context
	cancel context.CancelFunc
}

func (j join) Get() user.Token {
	return j.Token
}

func (j join) Inbox() party.PostChan {
	return j.inbox
}

func (j join) Cancel() context.CancelFunc {
	return j.cancel
}

func (j join) Reply() chan party.JoinReply {
	return j.reply
}

type leave struct {
	user.Token
}

func (l leave) Get() user.Token {
	return l.Token
}
