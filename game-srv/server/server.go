// Package server is related to the Cobra cmd serve. It
// contains a basic HTTP and WS server
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

// http config
const (
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
	maxHeaderBytes = 1 << 20
)

// Config of serve cmd
type Config struct {
	Version string
	Port    int
}

// Start the application
func Start(c Config) {
	slog.Info("server started",
		slog.String("version", c.Version),
		slog.Int("port", c.Port),
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	var ready atomic.Bool // TODO better logic

	s := serve(ctx, c.Port, &ready)

	ready.Store(true) // setup done

	// waits for app to interrupt
	<-ctx.Done()
	gracefulShutdown(s)
}

func gracefulShutdown(s *http.Server) {
	slog.Info("received interrupt")

	ctxTo, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// > program doesn't exit and waits instead for Shutdown to return.
	if err := s.Shutdown(ctxTo); err != nil {
		slog.Error("http shutdown", "error", err)
	}

	slog.Info("server terminated")
}

// serve the http server
func serve(ctx context.Context, port int, ready *atomic.Bool) *http.Server {
	handler, notifier := setupRoutes(ctx, ready)

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        handler,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
		// uses the default slog as log
		ErrorLog: slog.NewLogLogger(
			slog.Default().Handler(), slog.LevelInfo),
		BaseContext: func(l net.Listener) context.Context { return ctx },
	}

	// 🚀
	go func() {
		defer notifier.Close()
		slog.Error("listen and serve", "error", s.ListenAndServe())
	}()

	return s
}

func setupRoutes(
	ctx context.Context,
	ready *atomic.Bool,
) (http.Handler, ws.Notifier) {
	mux := http.NewServeMux()

	notifier := ws.New(ctx)

	mux.HandleFunc("GET /healthz", healthHandler())
	mux.HandleFunc("GET /readyz", readyHandler(ready))
	mux.HandleFunc("GET /ws", notifier.Handler())

	stack := middleware.Stack(
		middleware.Log,
	)

	return stack(mux), notifier
}
