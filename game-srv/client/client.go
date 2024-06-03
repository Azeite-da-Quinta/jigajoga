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
	Version   string
	Host      string
	NbWorkers int
	NbWrites  int
}

// Dial
func Dial(conf Config) {
	ctxS, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ctx, cancel := context.WithTimeout(ctxS, 1*time.Minute)
	defer cancel()

	if err := httpReady(conf); err != nil {
		return
	}

	var wg sync.WaitGroup
	for i := range conf.NbWorkers {
		wg.Add(1)
		go runWS(ctx, i,
			conf.NbWrites, &wg,
			fmt.Sprintf("ws://%s/ws", conf.Host),
		)
	}

	wg.Wait()
	slog.Info("closing clients")
	time.Sleep(100 * time.Millisecond)
}

func runWS(ctx context.Context, id, nbWrites int, wg *sync.WaitGroup, host string) {
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

	go write(ctx, ws, id, nbWrites)
	read(ctx, ws, id)
	slog.Info("work done", "worker", id)
}

func read(ctx context.Context, ws *websocket.Conn, id int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		ws.SetReadDeadline(time.Now().Add(10 * time.Second))

		mt, b, err := ws.ReadMessage()
		if err != nil {
			slog.Error("ws reading message", "id", id, "error", err)
			return
		}

		slog.Info("ws reading", "id", id, "type", mt, "msg", string(b))
	}
}

func write(ctx context.Context, ws *websocket.Conn, id, nbWrites int) {
	for i := range nbWrites {
		select {
		case <-ctx.Done():
			return
		default:
		}

		ws.SetWriteDeadline(time.Now().Add(10 * time.Second))

		err := ws.WriteMessage(
			websocket.TextMessage,
			[]byte(fmt.Sprintf("olá I'm worker %d sending %d", id, i)))
		if err != nil {
			slog.Error("failed to write", "id", id, "error", err)
		}

		// time.Sleep(1 * time.Second)
	}

	slog.Info("worker done writing", "id", id)
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
