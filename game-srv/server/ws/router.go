package ws

import (
	"context"
	"fmt"
	"log/slog"
)

type requestKind int

const (
	join requestKind = iota
	leave
)

type request struct {
	client *client
	kind   requestKind
}

// Router routes messages according to rooms
type Router struct {
	clients  map[string]*client
	messages chan []byte
	requests chan request
}

func NewRouter() *Router {
	return &Router{
		clients:  make(map[string]*client),
		messages: make(chan []byte, 10),  // TODO unbuffered ?
		requests: make(chan request, 10), // TODO unbuffered ?
	}
}

func (r *Router) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-r.requests:
			r.handleReq(req)
		}
	}
}

func (r *Router) handleReq(req request) {
	fmt.Println("handleReq", req)
	switch req.kind {
	case join:
		r.clients[req.client.id] = req.client

		slog.Info("client joined", "name", req.client.name)
		r.broadcast([]byte(fmt.Sprintf("%s joined", req.client.name)))
	case leave:
		delete(r.clients, req.client.id)

		slog.Info("client left", "name", req.client.name)
		r.broadcast([]byte(fmt.Sprintf("%s left", req.client.name)))
	default:
		fmt.Println("rien")
	}
}

func (r *Router) broadcast(b []byte) {
	for k, cl := range r.clients {
		select {
		case cl.inbox <- b:
		default:
			close(cl.inbox)
			delete(r.clients, k)
		}
	}
}
