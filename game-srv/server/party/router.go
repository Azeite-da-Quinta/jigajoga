package party

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/server/user"
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
			rttl.cancel()
			close(rttl.roomChan)
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

const (
	ttl = 3 * time.Hour
)

func (rt *Router) handleReq(ctx context.Context, req Request) {
	switch v := req.(type) {
	case Join:
		if rttl, ok := rt.rooms[v.RoomID()]; !ok {
			slog.Info("room created", "id", v.RoomID())

			ctxRoom, cancel := context.WithCancel(ctx)
			ch := make(chan Request)

			rt.rooms[v.RoomID()] = ttlRoom{
				roomChan:  ch,
				createdAt: time.Now(),
				cancel:    cancel,
			}

			r := newRoom(v.RoomID(), ch)
			go r.run(ctxRoom)

			ch <- req // forward to room
		} else {
			rttl.roomChan <- req // forward to room
		}
	case Leave:
		rttl, ok := rt.rooms[v.RoomID()]
		if ok {
			rttl.roomChan <- req // forward to room
		}
	default:
		slog.Warn("router handler: request type not handled")
	}
}

func (rt *Router) checkTTL(t time.Time) {
	for key, rttl := range rt.rooms {
		if t.Before(rttl.createdAt.Add(ttl)) {
			slog.Debug("room closing. ttl")
			rttl.cancel()
			close(rttl.roomChan)
			delete(rt.rooms, key)
		}
	}
}

// ttlRoom is a wrapper around a running room
type ttlRoom struct {
	createdAt time.Time
	// write only channel to running room
	roomChan chan<- Request
	cancel   context.CancelFunc
}
