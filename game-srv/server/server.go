package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/server/middleware"
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/server/ws"
)

type Config struct {
	Version string
	Port    int
}

func Start(c Config) {
	slog.Info("server started",
		slog.String("version", c.Version),
		slog.Int("port", c.Port),
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	var ready atomic.Bool
	s := Serve(ctx, c.Port, &ready)

	// setup done
	ready.Store(true)

	// waits for app to interrupt
	<-ctx.Done()
	slog.Info("closing: received interrupt")

	ctxTo, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// > program doesn't exit and waits instead for Shutdown to return.
	if err := s.Shutdown(ctxTo); err != nil {
		slog.Error("http shutdown", "error", err)
	}

	slog.Info("closing: server terminated")
}

func Serve(ctx context.Context, port int, ready *atomic.Bool) *http.Server {
	mux := http.NewServeMux()

	stack := middleware.Stack(
		middleware.Log,
	)

	notifier := ws.New(ctx)

	mux.HandleFunc("GET /healthz", healthHandler())
	mux.HandleFunc("GET /readyz", readyHandler(ready))
	mux.HandleFunc("GET /ws", notifier.Handler())

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        stack(mux),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		// uses the default slog as log
		ErrorLog:    slog.NewLogLogger(slog.Default().Handler(), slog.LevelInfo),
		BaseContext: func(l net.Listener) context.Context { return ctx },
	}

	// 🚀
	go func() {
		defer notifier.Close()
		slog.Error("listen and serve", "error", s.ListenAndServe())
	}()

	return s
}
