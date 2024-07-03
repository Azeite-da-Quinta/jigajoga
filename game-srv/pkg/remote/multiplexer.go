// Package remote implements a multiplexer of clients.
// From the perspective of the room, it will be like interracting with
// many, but in truth they will all be consumed from here.
package remote

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/party/comms"
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/party/event"
	"github.com/Azeite-da-Quinta/jigajoga/libs/chanx"
	"github.com/Azeite-da-Quinta/jigajoga/libs/envelope"
	"github.com/Azeite-da-Quinta/jigajoga/libs/slogt"
	"github.com/Azeite-da-Quinta/jigajoga/libs/user"
)

// Payload coming from remote
type Payload interface {
	Type() string
	User() user.Token
}

// Remote allows to subscribe to events
type Remote interface {
	Publish(ctx context.Context, e event.Envelope)
	SubRooms(context.Context) <-chan Payload
	SubMessages(ctx context.Context) <-chan envelope.Message
}

// Notifier allows to forward notifications on the chain
type Notifier interface {
	Forward(e string, t user.Token) (<-chan event.Reply, error)
}

// Multiplexer acts as an intermediate between rooms and remote clients
type Multiplexer struct {
	// receive replies when joining a room
	replies chan event.Reply
	// receives messages from various rooms
	multiInbox chan comms.Message
	// events coming from local emitters
	localEvents <-chan event.Query
	// Used to pass around notifications from Remote
	Notifier
}

const (
	repliesBuf = 1024
	multiBuf   = 8192
)

// New initializes a Multiplexer
func New(
	n Notifier,
	events <-chan event.Query, // from Tee
) Multiplexer {
	return Multiplexer{
		multiInbox:  make(chan comms.Message, multiBuf),
		replies:     make(chan event.Reply, repliesBuf),
		localEvents: events,
		Notifier:    n,
	}
}

const workUnit = 1

// Run the multiplexer work loop
func (mx *Multiplexer) Run(ctx context.Context, remote Remote) {
	// protect channels from being closed before writers
	var wg sync.WaitGroup
	defer wg.Wait()

	subMessages := remote.SubMessages(ctx)

	go mx.roomLoop(ctx, remote)

	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "mux: closing")
			return
		case m := <-subMessages:
			fmt.Println("m", m)
		case r := <-mx.replies:
			//
			mx.handleReply(ctx, r, &wg)
		case m := <-mx.multiInbox:
			fmt.Println("multiinbox", m)
			// remote <- rooms
			remote.Publish(ctx, m.Envelope())
		}
	}
}

func (mx *Multiplexer) roomLoop(ctx context.Context, remote Remote) {
	var wg sync.WaitGroup
	defer wg.Wait()

	subRooms := remote.SubRooms(ctx)

	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "mux: closing room loop")
			return
		case e := <-mx.localEvents:
			// remote <- router
			remote.Publish(ctx, e)
		case p := <-subRooms:
			p.User().ID() // TODO?
			// router <- remote
			mx.forwardEvent(ctx, p, &wg)
		}
	}

}

// forward event to Notifier
func (mx *Multiplexer) forwardEvent(
	ctx context.Context,
	p Payload,
	wg *sync.WaitGroup,
) {
	reply, err := mx.Forward(p.Type(), p.User())
	if err != nil {
		slog.Error("mux: failed to forward", slogt.Error(err))
		return
	}

	if reply != nil {
		wg.Add(workUnit)
		go chanx.Drain(
			ctx,
			reply,
			wg,
			mx.replies,
		)
	}
}

func (mx *Multiplexer) handleReply(
	ctx context.Context,
	r event.Reply,
	wg *sync.WaitGroup,
) {
	wg.Add(workUnit)

	go chanx.Drain(
		ctx,
		r.Inbox(),
		wg,
		mx.multiInbox,
	)
}
