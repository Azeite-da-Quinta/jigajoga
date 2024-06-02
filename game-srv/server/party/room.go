package party

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/server/user"
)

// room handles the multicast of messages between clients
type room struct {
	id user.Identifier
	// map by id of all client's write channels
	clients map[user.Identifier]PostChan
	// read incoming requests
	requests <-chan Request
	// read incoming messages to multicast to the room
	multicast chan []byte
}

// newRoom creates a new room with the required fields
func newRoom(id user.Identifier, c <-chan Request) room {
	return room{
		id:        id,
		clients:   make(map[user.Identifier]PostChan),
		requests:  c,
		multicast: make(chan []byte),
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

		r.clients[v.ID()] = v.ClientInbox

		// reply with the channel this room is listening on
		v.ReplyRoom <- r.multicast

		r.sendToAll([]byte(fmt.Sprintf("%s joined", v.Name())))
	case Leave:
		slog.Debug("client left",
			"id", v.ID(),
			"room", v.RoomID(),
			"name", v.Name())

		// closing emit ch
		if ch, ok := r.clients[v.ID()]; ok {
			close(ch)
		}

		delete(r.clients, v.ID())

		r.sendToAll([]byte(fmt.Sprintf("%s left", v.Name())))
	default:
		slog.Debug("unkown request kind in room handler")
	}
}

// sendToAll bytes message
func (r *room) sendToAll(b []byte) {
	for id, ch := range r.clients {
		select {
		case ch <- b:
		default:
			// if unavailable
			// 🔬 it will close if the client is unavailable
			close(ch)
			delete(r.clients, id)
		}
	}
}

// closeAll remaining client-write channels
func (r *room) closeAll() {
	for _, ch := range r.clients {
		close(ch)
	}

	clear(r.clients)
}
