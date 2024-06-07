// Package party implements the logic to route clients into Rooms
package party

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/server/user"
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
	requests <-chan Request
}

// NewRouter creates a New party Router
func NewRouter(c <-chan Request) Router {
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
	fmt.Println("calling close all")
	for _, rttl := range rt.rooms {
		go rttl.close()
	}

	clear(rt.rooms)
}

func (rt *Router) handleReq(ctx context.Context, req Request) {
	switch val := req.(type) {
	case Join:
		rt.handleJoin(ctx, val)
	case Leave:
		rt.handleLeave(val)
	default:
		slog.Warn("router handler: request type not handled")
	}
}

func (rt *Router) handleJoin(ctx context.Context, req Join) {
	if room, ok := rt.rooms[req.RoomID()]; !ok {
		slog.Info("room created", "id", req.RoomID())

		tr := makeRoom(ctx, req)
		rt.rooms[req.RoomID()] = tr

		forward(req, tr)
	} else {
		forward(req, room)
	}
}

func (rt *Router) handleLeave(req Leave) {
	rttl, ok := rt.rooms[req.RoomID()]
	if ok {
		rttl.requests <- req // forward to room
	}
}

// forward join request
func forward(req Join, tr ttlRoom) {
	//revive:disable:add-constant
	tr.wg.Add(1)

	tr.requests <- req // forward to room

	// reply with the channel this room is listening on
	req.Reply() <- JoinReply{
		RoomInbox: tr.messages,
		Wg:        tr.wg,
	}
}

func (rt *Router) checkTTL(t time.Time) {
	for key, rttl := range rt.rooms {
		if t.Before(rttl.createdAt.Add(ttl)) {
			// TODO log message with ID
			slog.Info("room closing. ttl")

			go rttl.close()

			delete(rt.rooms, key)
		}
	}
}

// makeRoom creates and runs a new room. Returns the wrapper
// that controls cancelation and channels
func makeRoom(ctx context.Context, jreq Join) ttlRoom {
	ctxRoom, cancel := context.WithCancel(ctx)
	requests := make(chan Request, roomBufSize)
	messages := make(chan []byte, messagesBufSize)

	tr := ttlRoom{
		requests:  requests,
		createdAt: time.Now(),
		cancel:    cancel,
		messages:  messages,
		wg:        &sync.WaitGroup{},
	}

	r := newRoom(jreq.RoomID(), roomChans{
		messages: messages,
		requests: requests,
	})
	go r.run(ctxRoom)

	return tr
}

// ttlRoom is a wrapper around a running room
type ttlRoom struct {
	createdAt time.Time
	// write only channel to submit requests to running room
	requests chan<- Request
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
	close(tr.messages)
}
