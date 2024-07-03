package client

import (
	"bytes"
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/party/comms"
	"github.com/Azeite-da-Quinta/jigajoga/libs/envelope"
	"github.com/Azeite-da-Quinta/jigajoga/libs/slogt"
	"github.com/gorilla/websocket"
)

func (p *IOPump) readPump(
	ctx context.Context,
	conn *websocket.Conn,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	setLimits(conn)

	for {
		b, err := read(conn)
		if err != nil {
			return
		}

		msg, err := parse(b)
		if err != nil {
			slog.Error("ws: failed to read envelope from bytes",
				slogt.Error(err))
			continue
		}
		msg.Sender = int64(p.ID())

		// TODO the default doesn't convince me here
		select {
		case <-ctx.Done():
			slog.Info("ws: read pump done")
			return
		case p.room <- msg:
		default:
		}
	}
}

func read(conn *websocket.Conn) ([]byte, error) {
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

func parse(b []byte) (comms.Message, error) {
	e, err := envelope.FromBytes(b)
	if err != nil {
		return comms.Message{}, err
	}

	return comms.FromEnvelope(e)
}
