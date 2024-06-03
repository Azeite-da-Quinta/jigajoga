package party

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/server/user"
)

const (
	ttl = 3 * time.Hour
)

// read only
type InboxChan <-chan []byte

// write only
type PostChan chan<- []byte

// Router routes messages according to rooms
type Router struct {
	rooms    map[user.Identifier]ttlRoom
	requests <-chan Request
}

func NewRouter(c <-chan Request) Router {
	return Router{
		rooms:    make(map[user.Identifier]ttlRoom),
		requests: c,
	}
}

// Run the Router's loop
func (rt *Router) Run(ctx context.Context) {
	// cancel all rooms and close all room writing channels
	defer func() {
		for _, rttl := range rt.rooms {
			rttl.close()
		}

		clear(rt.rooms)
	}()

	t := time.NewTicker(10 * time.Minute)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case req := <-rt.requests:
			rt.handleReq(ctx, req)
		case t := <-t.C:
			fmt.Println("check ttl")
			rt.checkTTL(t)
		}
	}
}

func (rt *Router) handleReq(ctx context.Context, req Request) {
	switch v := req.(type) {
	case Join:
		if rttl, ok := rt.rooms[v.RoomID()]; !ok {
			slog.Info("room created", "id", v.RoomID())

			tr := rt.makeRoom(ctx, v)
			rt.rooms[v.RoomID()] = tr

			// reply with the channel this room is listening on
			v.ReplyRoom <- JoinReply{
				RoomInbox: tr.messages,
			}

			tr.requests <- req // forward to room
		} else {
			// reply with the channel this room is listening on
			v.ReplyRoom <- JoinReply{
				RoomInbox: rttl.messages,
			}

			rttl.requests <- req // forward to room
		}
	case Leave:
		rttl, ok := rt.rooms[v.RoomID()]
		if ok {
			rttl.requests <- req // forward to room
		}
	default:
		slog.Warn("router handler: request type not handled")
	}
}

func (rt *Router) checkTTL(t time.Time) {
	for key, rttl := range rt.rooms {
		if t.Before(rttl.createdAt.Add(ttl)) {
			slog.Debug("room closing. ttl")

			rttl.close()

			delete(rt.rooms, key)
		}
	}
}

// makeRoom creates and runs a new room. Returns the wrapper
// that controls cancelation and channels
func (rt *Router) makeRoom(ctx context.Context, jreq Join) ttlRoom {
	ctxRoom, cancel := context.WithCancel(ctx)
	requests := make(chan Request)
	messages := make(chan []byte)

	tr := ttlRoom{
		requests:  requests,
		createdAt: time.Now(),
		cancel:    cancel,
		messages:  messages,
	}

	r := newRoom(jreq.RoomID(), requests, messages)
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
}

func (tr ttlRoom) close() {
	tr.cancel()
	close(tr.requests)
	// Validate clients cannot write here any longer
	close(tr.messages)
}
