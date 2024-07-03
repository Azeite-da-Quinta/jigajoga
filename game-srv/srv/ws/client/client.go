// Package client specific behavior for WS
package client

import (
	"context"
	"sync"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/party/event"
	"github.com/Azeite-da-Quinta/jigajoga/libs/user"
	"github.com/gorilla/websocket"
)

// read settings
const (
	pongWait       = 60 * time.Second
	maxMessageSize = 512
)

const (
	// clientBuf 🔬 controls how many messages can fill
	// the client's inbox before closing
	clientBuf = 16
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// IOPump writes and reads to conn
type IOPump struct {
	user.Token
	inbox event.InboxChan
	room  event.PostChan
}

// New creates a new IOPump
func New(t user.Token) IOPump {
	return IOPump{
		Token: t,
	}
}

// Pump starts the two pump workers
func (p *IOPump) Pump(
	ctx context.Context,
	conn *websocket.Conn,
	reply <-chan event.Reply,
) {
	r := <-reply

	p.inbox = r.Inbox()
	p.room = r.Room()

	// read and write
	const numWorkers = 2

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	go p.writePump(ctx, conn, &wg)
	go p.readPump(ctx, conn, &wg)

	wg.Wait()
}
