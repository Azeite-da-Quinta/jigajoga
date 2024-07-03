package room

import (
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/party/event"
)

type client struct {
	inbox event.PostChan
	// cancel context.CancelFunc
}

// close client
func (cl client) close() {
	// cl.cancel()
	close(cl.inbox)
}
