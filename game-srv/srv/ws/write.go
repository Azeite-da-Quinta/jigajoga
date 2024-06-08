package ws

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/libs/slogt"
	"github.com/gorilla/websocket"
)

func (cl *client) writePump(
	ctx context.Context,
	conn *websocket.Conn,
	wg *sync.WaitGroup,
) {
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
				slog.Debug("ws: inbox was closed by writer", "client", cl.ID())
				emitControl(conn, websocket.CloseMessage)
				return
			}

			if err := cl.handleInbox(conn, msg); err != nil {
				slog.Error("ws: return from write pump", slogt.Error(err))
				return
			}
		case <-pingT.C:
			if err := emitControl(conn, websocket.PingMessage); err != nil {
				return
			}
		}
	}
}

func (cl *client) handleInbox(conn *websocket.Conn, msg []byte) error {
	conn.SetWriteDeadline(time.Now().Add(writeDelay))
	err := conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return err
	}

	return cl.sendQueued(conn)
}

// sendQueued outputs all queued remaining messages
func (cl *client) sendQueued(conn *websocket.Conn) error {
	for range len(cl.inbox) {
		err := conn.WriteMessage(websocket.TextMessage, <-cl.inbox)
		if err != nil {
			return err
		}
	}

	return nil
}

func emitControl(conn *websocket.Conn, messageType int) error {
	return conn.WriteControl(
		messageType,
		nil,
		time.Now().Add(writeDelay),
	)
}
