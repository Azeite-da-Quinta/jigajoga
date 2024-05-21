package client

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type Config struct {
	Version string
	Host    string
}

func Dial(conf Config) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	c := http.Client{Timeout: time.Duration(1) * time.Second}

	resp, err := c.Get(fmt.Sprintf("http://%s/readyz", conf.Host))
	if err != nil {
		slog.Error("failed to http get ready", "error", err)
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to check ready", "error", err)
		return
	}
	slog.Info("http response", "status", resp.Status, "bytes", string(b))

	ws, respws, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws", conf.Host), nil)
	if err != nil {
		slog.Error("failed to dial ws", "error", err)
		return
	}
	defer ws.Close()

	ws.SetCloseHandler(func(code int, text string) error {
		slog.Info("ws connection closed", "code", code, "reason", text)
		return nil
	})

	slog.Info("ws response", "status", respws.Status)

	for range 10 {
		select {
		case <-ctx.Done():
			return
		default:
		}

		ws.WriteMessage(websocket.TextMessage, []byte("olá"))
		mt, b, err := ws.ReadMessage()
		if err != nil {
			slog.Error("ws reading message", "error", err)
			return
		}

		slog.Info("ws reading", "type", mt, "msg", string(b))
		time.Sleep(2 * time.Second)
	}

	<-ctx.Done()
}
