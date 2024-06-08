package ws

import (
	"context"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/srv/party"
	"github.com/Azeite-da-Quinta/jigajoga/libs/token"
	"github.com/Azeite-da-Quinta/jigajoga/libs/user"
	"github.com/gorilla/websocket"
)

const (
	// clientBuf 🔬 controls how many messages can fill
	// the client's inbox before closing
	clientBuf = 16
)

// New returns a Notifier and runs the party Router
func New(ctx context.Context, secretKey string) (Notifier, error) {
	c := make(chan party.Request, party.RouterBufSize)

	b, err := token.Base64ToKey(secretKey)
	if err != nil {
		return Notifier{}, err
	}

	n := Notifier{
		requests: c,
		codec:    token.Codec{Key: b},
	}

	rt := party.NewRouter(c)
	go rt.Run(ctx)

	return n, nil
}

// Notifier is responsible to link the websocket connection
// to rooms
type Notifier struct {
	// Submit requests to Router
	requests chan<- party.Request
	codec    token.Codec
}

// join a client to a room together.
func (n *Notifier) join(ctx context.Context,
	conn *websocket.Conn,
	t user.Token) {
	ch := make(chan []byte, clientBuf)
	//revive:disable:add-constant
	reply := make(chan party.JoinReply, 1)

	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	n.requests <- join{
		Token:  t,
		inbox:  ch,
		reply:  reply,
		cancel: cancel,
	}

	resp := <-reply

	defer func() {
		resp.Wg.Done()
		n.requests <- leave{
			Token: t,
		}
	}()

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
