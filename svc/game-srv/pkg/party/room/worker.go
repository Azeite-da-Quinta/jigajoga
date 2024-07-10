// Package room contains the code to multicast between clients
package room

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Azeite-da-Quinta/jigajoga/libs/go/slogt"
	"github.com/Azeite-da-Quinta/jigajoga/libs/go/user"
	"github.com/Azeite-da-Quinta/jigajoga/svc/game-srv/pkg/party/comms"
	"github.com/Azeite-da-Quinta/jigajoga/svc/game-srv/pkg/party/event"
)

// Worker handles the multicast of messages between clients
type Worker struct {
	// map by id of all client's write channels
	clients map[user.Identifier]client
	// read incoming requests
	requests <-chan event.Query
	// read incoming messages from clients to the room
	multicast chan comms.Message
	id        user.Identifier
}

// Parameters to configure a room
type Parameters struct {
	Requests <-chan event.Query
}

const messagesBufSize = 512

// New creates a new room worker
func New(
	id user.Identifier,
	requests <-chan event.Query,
) Worker {
	return Worker{
		id:        id,
		clients:   make(map[user.Identifier]client),
		requests:  requests,
		multicast: make(chan comms.Message, messagesBufSize),
	}
}

// Run the room's loop
func (w *Worker) Run(ctx context.Context) {
	defer func() {
		slog.Debug("room: closing",
			slogt.Room(int64(w.id)))

		w.sendToAll(comms.Message{
			Content: "room closed",
		})
		w.closeAll()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case req := <-w.requests:
			w.handleReq(req)
		case b := <-w.multicast:
			/* slog.Debug("room: received bytes from multicast",
			"content", b.Content, slogt.PlayerID(b.Sender)) */
			w.sendToAll(b)
		}
	}
}

const (
	// clientBuf 🔬 controls how many messages can fill
	// the client's inbox before closing
	clientBuf = 16
)

func (w *Worker) handleReq(evt event.Query) {
	switch e := evt.(type) {
	case event.Join:
		// cancel: e.Cancel(),
		w.handleJoin(e)
	case event.Leave:
		// TODO what are those comments already?
		// closing emit ch
		w.handleLeave(e)
	default:
		slog.Debug("room: unkown request kind in room handler")
	}
}

func (w *Worker) handleJoin(e event.Join) {
	defer close(e.Reply())

	slog.Debug("room: client joined",
		slogt.PlayerID(int64(e.ID())),
		slogt.Room(int64(e.RoomID())),
		slogt.PlayerName(e.Name()))

	ch := make(chan comms.Message, clientBuf)

	w.clients[e.ID()] = client{
		inbox: ch,
	}

	e.Reply() <- reply{
		id:    int64(w.id),
		post:  w.multicast,
		inbox: ch,
	}

	w.sendToAll(
		comms.Message{
			// TODO complete
			Content: fmt.Sprintf("%s joined", e.Name()),
		},
	)
}

func (w *Worker) handleLeave(e event.Leave) {
	slog.Debug("room: client left",
		slogt.PlayerID(int64(e.ID())),
		slogt.Room(int64(e.RoomID())),
		slogt.PlayerName(e.Name()))

	if cl, ok := w.clients[e.ID()]; ok {
		cl.close()
	}
	delete(w.clients, e.ID())

	w.sendToAll(
		comms.Message{
			Content: fmt.Sprintf("%s left", e.Name()),
		},
	)
}

// sendToAll bytes message
func (w *Worker) sendToAll(msg comms.Message) {
	for id, cl := range w.clients {
		select {
		case cl.inbox <- msg:
		default:
			// if unavailable
			// 🔬 it will close if the client is unavailable
			cl.close()
			delete(w.clients, id)
		}
	}
}

// closeAll remaining client-write channels
func (w *Worker) closeAll() {
	for _, cl := range w.clients {
		cl.close()
	}

	clear(w.clients)
}
