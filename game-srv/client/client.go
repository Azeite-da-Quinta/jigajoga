package client

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Config struct {
	Version string
	Host    string
}

// Dial
func Dial(conf Config) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := httpReady(conf); err != nil {
		return
	}

	var wg sync.WaitGroup
	for i := range 10 {
		wg.Add(1)
		go runWS(ctx, i, &wg, fmt.Sprintf("ws://%s/ws", conf.Host))
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond)
}

func runWS(ctx context.Context, _ int, wg *sync.WaitGroup, host string) {
	defer wg.Done()

	ws, respws, err := websocket.DefaultDialer.Dial(host, http.Header{})
	if err != nil {
		slog.Error("failed to dial ws", "error", err)
		return
	}
	defer ws.Close()

	slog.Info("ws response", "status", respws.Status)

	ws.SetCloseHandler(func(code int, text string) error {
		slog.Info("ws connection closed", "code", code, "reason", text)
		return nil
	})

	go write(ctx, ws)
	read(ctx, ws)
}

func read(ctx context.Context, ws *websocket.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		mt, b, err := ws.ReadMessage()
		if err != nil {
			slog.Error("ws reading message", "error", err)
			return
		}
		slog.Info("ws reading", "type", mt, "msg", string(b))
	}
}

func write(ctx context.Context, ws *websocket.Conn) {
	for range 10 {
		select {
		case <-ctx.Done():
			return
		default:
		}

		ws.WriteMessage(websocket.TextMessage, []byte("olá"))

		time.Sleep(1 * time.Second)
	}
}

func httpReady(conf Config) error {
	c := http.Client{Timeout: time.Duration(1) * time.Second}

	resp, err := c.Get(fmt.Sprintf("http://%s/readyz", conf.Host))
	if err != nil {
		slog.Error("failed to http get ready", "error", err)
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to check ready", "error", err)
		return err
	}
	slog.Info("http response", "status", resp.Status, "bytes", string(b))

	return nil
}
