package party

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/server/user"
)

// room handles the multicast of messages between clients
type room struct {
	// map by id of all client's write channels
	clients map[user.Identifier]client
	// read incoming requests
	requests <-chan Request
	// read incoming messages from clients to the room
	multicast <-chan []byte
	id        user.Identifier
}

type roomChans struct {
	requests <-chan Request
	messages <-chan []byte
}

// newRoom creates a new room with the required fields
func newRoom(id user.Identifier, c roomChans) room {
	return room{
		id:        id,
		clients:   make(map[user.Identifier]client),
		requests:  c.requests,
		multicast: c.messages,
	}
}

// run the room's loop
func (r *room) run(ctx context.Context) {
	defer func() {
		slog.Debug("room closing",
			"room", r.id)

		r.sendToAll([]byte("room closed"))
		r.closeAll()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case req := <-r.requests:
			r.handleReq(req)
		case b := <-r.multicast:
			slog.Debug("received bytes from multicast", "bytes", string(b))
			r.sendToAll(b)
		}
	}
}

func (r *room) handleReq(req Request) {
	switch v := req.(type) {
	case Join:
		slog.Debug("client joined",
			"id", v.ID(),
			"room", v.RoomID(),
			"name", v.Name())

		r.clients[v.ID()] = client{
			inbox:  v.Inbox(),
			cancel: v.Cancel(),
		}

		r.sendToAll([]byte(fmt.Sprintf("%s joined", v.Name())))
	case Leave:
		slog.Debug("client left",
			"id", v.ID(),
			"room", v.RoomID(),
			"name", v.Name())

		// closing emit ch
		if cl, ok := r.clients[v.ID()]; ok {
			cl.close()
		}
		delete(r.clients, v.ID())

		r.sendToAll([]byte(fmt.Sprintf("%s left", v.Name())))
	default:
		slog.Debug("unkown request kind in room handler")
	}
}

// sendToAll bytes message
func (r *room) sendToAll(b []byte) {
	for id, cl := range r.clients {
		select {
		case cl.inbox <- b:
		default:
			// if unavailable
			// 🔬 it will close if the client is unavailable
			cl.close()
			delete(r.clients, id)
		}
	}
}

// closeAll remaining client-write channels
func (r *room) closeAll() {
	for _, cl := range r.clients {
		cl.close()
	}

	clear(r.clients)
}

type client struct {
	inbox  PostChan
	cancel context.CancelFunc
}

// close client
func (cl client) close() {
	cl.cancel()
	close(cl.inbox)
}
