package client

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/libs/slogt"
	"github.com/gorilla/websocket"
)

type worker struct {
	url, jwt      string
	num, nbWrites int
}

func (w worker) run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %v", w.jwt))

	ws, respws, err := websocket.DefaultDialer.Dial(w.url, header)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "failed to dial ws",
			slogt.Error(err), slog.String("status", respws.Status))
		return
	}
	defer ws.Close()

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

		ws.SetWriteDeadline(time.Now().Add(10 * time.Second))

		err := ws.WriteMessage(
			websocket.TextMessage,
			[]byte(fmt.Sprintf("olá I'm worker %d sending %d", w.num, i)))
		if err != nil {
			slog.Error("failed to write", slogt.Num(w.num), slogt.Error(err))
		}
	}

	slog.Info("worker done writing", slogt.Num(w.num))
}
