package ws

import (
	"bytes"
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/libs/slogt"
	"github.com/gorilla/websocket"
)

const (
	pongWait       = 60 * time.Second
	maxMessageSize = 512
)

func (cl *client) readPump(
	ctx context.Context,
	conn *websocket.Conn,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	setLimits(conn)

	for {
		message, err := readMessage(conn)
		if err != nil {
			return
		}

		select {
		case <-ctx.Done():
			return
		case cl.room <- message:
		default:
		}
	}
}

func readMessage(conn *websocket.Conn) ([]byte, error) {
	_, message, err := conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err,
			websocket.CloseGoingAway,
			websocket.CloseAbnormalClosure) {
			slog.Error("ws: read unexpected close", slogt.Error(err))
		}
		return nil, err
	}

	return bytes.TrimSpace(bytes.ReplaceAll(message, newline, space)), nil
}

func setLimits(conn *websocket.Conn) {
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}
