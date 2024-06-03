package ws

import (
	"context"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/server/party"
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/server/user"
	"github.com/gorilla/websocket"
)

// clientBuf 🔬 controls how many messages can fill
// the client's inbox before closing
const clientBuf = 10

func New(ctx context.Context) Notifier {
	c := make(chan party.Request, 1024)

	n := Notifier{
		requests: c,
	}

	rt := party.NewRouter(c)
	go rt.Run(ctx)

	return n
}

type Notifier struct {
	// Submit requests to Router
	requests chan<- party.Request
}

// join a client to a room together.
func (n *Notifier) join(ctx context.Context, conn *websocket.Conn, t user.Token) {
	ch := make(chan []byte, clientBuf)
	reply := make(chan party.JoinReply)

	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	n.requests <- party.Join{
		Token:       t,
		ClientInbox: ch,
		ReplyRoom:   reply,
		Cancel:      cancel,
	}

	defer func() {
		n.requests <- party.Leave{
			Token: t,
		}
	}()

	resp := <-reply

	c := client{
		Token: t,
		inbox: ch,
		room:  resp.RoomInbox,
	}

	c.pump(subCtx, conn)
}

// Close notifier resources
func (n *Notifier) Close() {
	close(n.requests)
}
