package notifier

import "github.com/Azeite-da-Quinta/jigajoga/svc/game-srv/pkg/party"

type reply struct {
	post  party.PostChan
	inbox party.InboxChan
}

func (r reply) Room() party.PostChan {
	return r.post
}

func (r reply) Inbox() party.InboxChan {
	return r.inbox
}
