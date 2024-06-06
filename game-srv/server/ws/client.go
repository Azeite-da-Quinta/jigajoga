package ws

import (
	"context"
	"sync"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/server/party"
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/server/user"
	"github.com/gorilla/websocket"
)

const (
	writeDelay     = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = pongWait - 6*time.Second // 90%
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type client struct {
	user.Token
	inbox party.InboxChan
	room  party.PostChan
}

// read and write
const wsWorkers = 2

func (cl *client) pump(ctx context.Context, conn *websocket.Conn) {
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(wsWorkers)
	go cl.writePump(ctx, conn, &wg)
	go cl.readPump(ctx, conn, &wg)

	wg.Wait()
}
