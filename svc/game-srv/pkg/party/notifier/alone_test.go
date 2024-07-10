package notifier

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/libs/go/user"
	"github.com/Azeite-da-Quinta/jigajoga/svc/game-srv/pkg/party"
	"github.com/Azeite-da-Quinta/jigajoga/svc/game-srv/pkg/party/comms"
	"github.com/Azeite-da-Quinta/jigajoga/svc/game-srv/pkg/party/event"
)

var (
	messages = genMessages(100)
)

// revive:disable:add-constant

func Test_OneClient(t *testing.T) {
	d, ok := t.Deadline()
	if !ok {
		d = time.Now().Add(30 * time.Second)
	}

	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	n, routed := New(ctx)

	rt := party.NewRouter(routed)
	go rt.Run(ctx)

	mc := mockClient{
		messages: messages,
		n:        n,
	}

	t.Run("one client", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(2)
		mc.n.Notify(
			user.MockToken(),
			func(reply <-chan event.Reply) {
				r := <-reply

				go mc.receiveLoop(ctx, &wg, r.Inbox())
				go mc.emitLoop(ctx, &wg, r.Room())
			},
		)
		wg.Wait()

		if d := mc.distributed.Load(); int(d) != len(mc.messages) {
			t.Errorf("distributed: %d", d)
		}
	})
}

type mockClient struct {
	n           *Notifier
	messages    []comms.Message
	distributed atomic.Int32
}

func (mc *mockClient) emitLoop(
	ctx context.Context,
	wg *sync.WaitGroup,
	room event.PostChan,
) {
	defer wg.Done()

	for _, msg := range mc.messages {
		select {
		case <-ctx.Done():
			return
		case room <- msg:
			mc.distributed.Add(1)
		}
	}
}

func (mc *mockClient) receiveLoop(
	ctx context.Context,
	wg *sync.WaitGroup,
	inbox event.InboxChan,
) {
	defer wg.Done()

	for range len(mc.messages) {
		select {
		case <-ctx.Done():
			return
		case <-inbox:
		}
	}
}

func genMessages(count int) []comms.Message {
	messages := make([]comms.Message, 0, count)
	for i := range count {
		strconv.Itoa(i)

		messages = append(messages,
			comms.Message{
				Content: strconv.Itoa(i),
			},
		)
	}

	return messages
}
