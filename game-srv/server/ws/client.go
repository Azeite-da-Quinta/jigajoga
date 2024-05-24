package ws

import (
	"bytes"
	"context"
	"log/slog"
	"sync"
	"time"

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
	inbox chan []byte
	id    string
	name  string
}

func (cl *client) run(ctx context.Context, rt *Router, conn *websocket.Conn) {
	slog.Debug("connection start")
	defer conn.Close()

	rt.requests <- request{client: cl, kind: join}

	defer func() {
		rt.requests <- request{client: cl, kind: leave}
	}()

	var wg sync.WaitGroup
	wg.Add(2)
	go cl.writePump(ctx, &wg, conn)
	go readPump(ctx, &wg, conn, rt)

	wg.Wait()
	slog.Debug("connection over")
}

func readPump(_ context.Context, wg *sync.WaitGroup, conn *websocket.Conn, rt *Router) {
	defer wg.Done()

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("ws: read unexpected close", "err", err)
			}
			return
		}
		message = bytes.TrimSpace(bytes.ReplaceAll(message, newline, space))
		rt.broadcast(message)
	}
}

func (cl *client) writePump(ctx context.Context, wg *sync.WaitGroup, conn *websocket.Conn) {
	defer wg.Done()

	pingT := time.NewTicker(pingPeriod)
	defer pingT.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case msg, ok := <-cl.inbox:
			// inbox closed
			if !ok {
				_ = emitControl(conn, websocket.CloseMessage)
				return
			}
			conn.SetWriteDeadline(time.Now().Add(writeDelay))

			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(msg)

			// queued remaining messages
			n := len(cl.inbox)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-cl.inbox)
			}

			// flush writer
			if err := w.Close(); err != nil {
				return
			}
		case <-pingT.C:
			if err := emitControl(conn, websocket.PingMessage); err != nil {
				return
			}
		}
	}
}

func emitControl(conn *websocket.Conn, messageType int) error {
	return conn.WriteControl(
		messageType,
		nil,
		time.Now().Add(writeDelay),
	)
}
