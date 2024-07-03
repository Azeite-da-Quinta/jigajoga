// Package notifier wraps how to notify router of queries
package notifier

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/party"
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/party/event"
	"github.com/Azeite-da-Quinta/jigajoga/libs/chanx"
	"github.com/Azeite-da-Quinta/jigajoga/libs/user"
)

// New returns a simple Notifier
func New(ctx context.Context) (*Notifier, <-chan event.Query) {
	c := make(chan event.Query, party.RouterBufSize)

	n := &Notifier{
		requests: c,
	}

	return n, chanx.OrDone(ctx, c)
}

// NewTee returns an enhanced Notifier
func NewTee(ctx context.Context) (n *Notifier, routed, spy <-chan event.Query) {
	direct := make(chan event.Query, party.RouterBufSize)
	requests := make(chan event.Query, party.RouterBufSize)

	n = &Notifier{
		requests: requests,
		direct:   direct,
	}

	/*
		routed notifications come from two sources (Fan-In):
		- requests, that are splited with a Tee so they can be duplicated
		- direct (skips)
	*/

	internal, spy := chanx.Tee(ctx, requests)
	routed = chanx.FanIn(ctx, direct, internal)

	return n, routed, spy
}

// Notifier is responsible to link the websocket connection
// to rooms
type Notifier struct {
	// Submit requests to Router
	requests chan<- event.Query
	// In Tee mode, this submits in direct.
	// DO NOT USE OTHERWISE.
	direct chan<- event.Query
	// to avoid notif duplication
	seen sync.Map
}

const replyBuf = 1

// Notify a client to a room together.
func (n *Notifier) Notify(
	t user.Token,
	fn func(reply <-chan event.Reply),
) {
	reply := make(chan event.Reply, replyBuf)

	n.seen.Store(t.ID(), empty)

	n.requests <- join{
		Token: t,
		reply: reply,
	}

	defer func() {
		n.seen.Delete(t.ID())

		n.requests <- leave{
			Token: t,
		}
	}()

	fn(reply)
}

// Forward a query from an arbitrary string event and user in Tee mode
func (n *Notifier) Forward(e string, t user.Token) (<-chan event.Reply, error) {
	_, ok := n.seen.Load(t.ID())
	if ok {
		return nil, nil
	}

	q, err := getQuery(e, t)
	if err != nil {
		return nil, err
	}

	var reply <-chan event.Reply
	if j, ok := q.(*join); ok {
		c := make(chan event.Reply, replyBuf)
		j.reply = c
		reply = c
	}

	n.direct <- q

	return reply, nil
}

// Close notifier resources
func (n *Notifier) Close() {
	close(n.requests)
}

func getQuery(s string, t user.Token) (event.Query, error) {
	switch strings.ToLower(s) {
	case "join":
		return &join{
			Token: t,
		}, nil
	case "leave":
		return &leave{
			Token: t,
		}, nil
	default:
		return join{}, errors.New("not implemented")
	}
}

var empty = struct{}{}
