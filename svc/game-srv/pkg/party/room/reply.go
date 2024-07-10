package room

import "github.com/Azeite-da-Quinta/jigajoga/svc/game-srv/pkg/party/event"

type reply struct {
	post  event.PostChan
	inbox event.InboxChan
	id    int64
}

func (r reply) OriginID() int64 {
	return r.id
}

func (r reply) Room() event.PostChan {
	return r.post
}

func (r reply) Inbox() event.InboxChan {
	return r.inbox
}
