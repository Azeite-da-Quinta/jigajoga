// Package party implements the logic to route clients into Rooms
package party

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/libs/go/slogt"
	"github.com/Azeite-da-Quinta/jigajoga/libs/go/user"
	"github.com/Azeite-da-Quinta/jigajoga/svc/game-srv/pkg/party/event"
	"github.com/Azeite-da-Quinta/jigajoga/svc/game-srv/pkg/party/room"
)

// RouterBufSize router's channel size (handles Join/leave requests)
const RouterBufSize = 1024

const (
	ttl             = 3 * time.Hour
	roomBufSize     = 32
	messagesBufSize = 512
)

// InboxChan is a read only message channel
type InboxChan <-chan []byte

// PostChan is a write only message channel
type PostChan chan<- []byte

// Router routes messages according to rooms
type Router struct {
	rooms    map[user.Identifier]ttlRoom
	requests <-chan event.Query
}

// NewRouter creates a New party Router
func NewRouter(c <-chan event.Query) Router {
	return Router{
		rooms:    make(map[user.Identifier]ttlRoom),
		requests: c,
	}
}

// Run the Router's loop
func (rt *Router) Run(ctx context.Context) {
	defer rt.closeAll()

	t := time.NewTicker(10 * time.Minute)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case req := <-rt.requests:
			rt.handleReq(ctx, req)
		case t := <-t.C:
			rt.checkTTL(t)
		}
	}
}

// closeAll cancel all rooms and close all room writing channels
func (rt *Router) closeAll() {
	for _, rttl := range rt.rooms {
		go rttl.close()
	}

	clear(rt.rooms)
}

func (rt *Router) handleReq(ctx context.Context, req event.Query) {
	switch val := req.(type) {
	case event.Join:
		rt.handleJoin(ctx, val)
	case event.Leave:
		rt.handleLeave(val)
	default:
		slog.Warn("router: request type not handled")
	}
}

func (rt *Router) handleJoin(ctx context.Context, req event.Join) {
	if rm, ok := rt.rooms[req.RoomID()]; !ok {
		slog.Info("router: room created", slogt.Room(int64(req.RoomID())))

		tr := makeRoom(ctx, req)
		rt.rooms[req.RoomID()] = tr

		forward(req, tr)
	} else {
		forward(req, rm)
	}
}

func (rt *Router) handleLeave(req event.Leave) {
	rttl, ok := rt.rooms[req.RoomID()]
	if ok {
		rttl.requests <- req // forward to room
	}
}

// forward join request
func forward(req event.Query, tr ttlRoom) {
	//revive:disable:add-constant
	tr.wg.Add(1)

	tr.requests <- req // forward to room
}

func (rt *Router) checkTTL(t time.Time) {
	for key, rttl := range rt.rooms {
		if t.Before(rttl.createdAt.Add(ttl)) {
			// TODO log message with ID
			slog.Info("router: room closing due to ttl")

			go rttl.close()

			delete(rt.rooms, key)
		}
	}
}

// makeRoom creates and runs a new room. Returns the wrapper
// that controls cancelation and channels
func makeRoom(ctx context.Context, jreq event.Join) ttlRoom {
	ctxRoom, cancel := context.WithCancel(ctx)
	requests := make(chan event.Query, roomBufSize)
	// messages := make(chan []byte, messagesBufSize)

	tr := ttlRoom{
		requests:  requests,
		createdAt: time.Now(),
		cancel:    cancel,
		// messages:  messages,
		wg: &sync.WaitGroup{},
	}

	r := room.New(jreq.RoomID(), requests)

	go r.Run(ctxRoom)

	return tr
}

// ttlRoom is a wrapper around a running room
type ttlRoom struct {
	createdAt time.Time
	// write only channel to submit requests to running room
	requests chan<- event.Query
	// cancel room's sub-ctx
	cancel context.CancelFunc
	// write only channel to submit messages to running room
	messages PostChan
	wg       *sync.WaitGroup
}

func (tr ttlRoom) close() {
	tr.cancel()
	close(tr.requests)
	tr.wg.Wait()
	// TODO Validate clients cannot write here any longer
	// close(tr.messages)
}
