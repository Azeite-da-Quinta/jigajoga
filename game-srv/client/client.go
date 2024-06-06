// Package client to test the server
package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/libs/slogt"
	"github.com/gorilla/websocket"
)

// Config of Dial
type Config struct {
	Version   string
	Host      string
	NbWorkers int
	NbWrites  int
}

// Dial connects to the server with N workers doing N jobs (write messages)
func Dial(conf Config) {
	ctxS, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ctx, cancel := context.WithTimeout(ctxS, 1*time.Minute)
	defer cancel()

	if err := httpReady(ctx, conf); err != nil {
		return
	}

	doJobs(ctx, conf)

	slog.Info("closing clients")
	time.Sleep(100 * time.Millisecond)
}

func doJobs(ctx context.Context, conf Config) {
	var wg sync.WaitGroup
	wg.Add(conf.NbWorkers)

	wconf := workerConf{
		url:      urlWS(conf.Host),
		nbWrites: conf.NbWrites,
	}

	for i := range conf.NbWorkers {
		wconf.id = i
		go runWS(ctx, wconf, &wg)
	}

	wg.Wait()
}

func urlWS(host string) string {
	return fmt.Sprintf("ws://%s/ws", host)
}

type workerConf struct {
	url          string
	id, nbWrites int
}

func runWS(ctx context.Context, wconf workerConf, wg *sync.WaitGroup) {
	defer wg.Done()

	ws, respws, err := websocket.DefaultDialer.Dial(wconf.url, http.Header{})
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

	go write(ctx, ws, wconf.id, wconf.nbWrites)
	read(ctx, ws, wconf.id)
	slog.Info("work done", "worker", wconf.id)
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
			slog.Error("failed to write", slogt.ID(id), slogt.Error(err))
		}
	}

	slog.Info("worker done writing", slogt.ID(id))
}

func urlReady(conf Config) string {
	return fmt.Sprintf("http://%s/readyz", conf.Host)
}

func httpReady(ctx context.Context, conf Config) error {
	c := http.Client{Timeout: time.Second}

	resp, err := c.Get(urlReady(conf))
	if err != nil {
		slog.Error("failed to http get ready", slogt.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.LogAttrs(ctx, slog.LevelError, "status not ok",
			slog.Int("status", resp.StatusCode))
		return errors.New("server not ready")
	}

	return nil
}
