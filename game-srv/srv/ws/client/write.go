package client

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/party/comms"
	"github.com/Azeite-da-Quinta/jigajoga/libs/slogt"
	"github.com/gorilla/websocket"
)

// write settings
const (
	writeDelay = 10 * time.Second
	pingPeriod = 15 * time.Second
)

func (p *IOPump) writePump(
	ctx context.Context,
	conn *websocket.Conn,
	wg *sync.WaitGroup,
) {
	defer conn.Close()
	defer wg.Done()

	pingT := time.NewTicker(pingPeriod)
	defer pingT.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pingT.C:
			if err := emitControl(conn, websocket.PingMessage); err != nil {
				return
			}
		case msg, ok := <-p.inbox:
			conn.SetWriteDeadline(time.Now().Add(writeDelay))

			// inbox closed
			if !ok {
				slog.Debug("ws: inbox was closed by writer", "client", p.ID())
				emitControl(conn, websocket.CloseMessage)
				return
			}

			b, err := marshal(msg)
			if err != nil {
				continue
			}

			if err := p.handleInbox(conn, b); err != nil {
				slog.Error("ws: return from write pump", slogt.Error(err))
				return
			}
		}
	}
}

func (p *IOPump) handleInbox(conn *websocket.Conn, msg []byte) error {
	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		return err
	}

	return p.sendQueued(conn)
}

// sendQueued outputs all queued remaining messages
func (p *IOPump) sendQueued(conn *websocket.Conn) error {
	for range len(p.inbox) {
		b, err := marshal(<-p.inbox)
		if err != nil {
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			return err
		}
	}

	return nil
}

func emitControl(conn *websocket.Conn, messageType int) error {
	err := conn.WriteControl(
		messageType,
		nil,
		time.Now().Add(writeDelay),
	)
	if err != nil {
		slog.Error("ws: failed to emit control", slogt.Error(err))
	}

	return err
}

func marshal(msg comms.Message) ([]byte, error) {
	b, err := msg.Envelope().Serialize()
	if err != nil {
		slog.Error("ws: failed to marshal", slogt.Error(err))
	}

	return b, err
}
