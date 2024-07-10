package client

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/libs/go/envelope"
	"github.com/Azeite-da-Quinta/jigajoga/libs/go/slogt"
	"github.com/gorilla/websocket"
)

type worker struct {
	url, jwt      string
	num, nbWrites int
}

func (w worker) run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	header := http.Header{}
	// Note: changed from Authorization header since it's not possible in regular web browsers
	header.Add("Sec-WebSocket-Protocol", fmt.Sprintf("base64url.bearer.authorization.%s, %s", w.jwt, envelope.Subprotocol))

	ws, respws, err := websocket.DefaultDialer.Dial(w.url, header)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "failed to dial ws",
			slogt.Error(err), slog.String("status", respws.Status))
		return
	}
	defer ws.Close()

	if respws.Header.Get("Sec-Websocket-Protocol") != envelope.Subprotocol {
		slog.Error("wrong ws subprotocol")
		return
	}

	ws.SetCloseHandler(func(code int, text string) error {
		slog.Info("ws connection closed", "code", code, "reason", text)
		return nil
	})

	go w.write(ctx, ws)
	w.read(ctx, ws)
	slog.Info("work done", "worker", w.num)
}

func (w worker) read(ctx context.Context, ws *websocket.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		ws.SetReadDeadline(time.Now().Add(10 * time.Second))

		mt, b, err := ws.ReadMessage()
		if err != nil {
			slog.Error("ws reading message", slogt.Num(w.num), slogt.Error(err))
			return
		}

		slog.Info("ws reading", slogt.Num(w.num), "type", mt, "msg", string(b))
	}
}

func (w worker) write(ctx context.Context, ws *websocket.Conn) {
	for i := range w.nbWrites {
		select {
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(200 * time.Millisecond)

		ws.SetWriteDeadline(time.Now().Add(10 * time.Second))

		e := envelope.Message{
			To:      "0",
			Content: fmt.Sprintf("olá I'm worker %d sending %d", w.num, i),
		}

		b, err := e.Serialize()
		if err != nil {
			slog.Error("failed to serialize",
				slogt.Num(w.num), slogt.Error(err))
			return
		}

		err = ws.WriteMessage(
			websocket.TextMessage,
			// TODO use the approriate envelope
			b,
		)
		if err != nil {
			slog.Error("failed to write", slogt.Num(w.num), slogt.Error(err))
		}
	}

	slog.Info("worker done writing", slogt.Num(w.num))
}
