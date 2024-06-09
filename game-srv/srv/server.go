// Package srv is related to the Cobra cmd serve. It
// contains a basic HTTP and WS server
package srv

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

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/srv/ws"
	"github.com/Azeite-da-Quinta/jigajoga/libs/handlers"
	"github.com/Azeite-da-Quinta/jigajoga/libs/middleware"
	"github.com/Azeite-da-Quinta/jigajoga/libs/slogt"
)

// http config
const (
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
	maxHeaderBytes = 1 << 20
)

// Config of serve cmd
type Config struct {
	Version   string
	JWTSecret string
	Port      int
}

// Server wraps an http server and more. It is setup by Config
type Server struct {
	httpSrv *http.Server
	Config
	ready atomic.Bool
}

// Start the application
func (s *Server) Start() {
	slog.Info("server started",
		slog.String("version", s.Version),
		slog.Int("port", s.Port),
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	err := s.serve(ctx)
	if err != nil {
		slog.Error("failed to serve http", slogt.Error(err))
		panic(err)
	}
	s.setReady()

	// waits for app to interrupt
	<-ctx.Done()
	s.gracefulShutdown()
}

func (s *Server) gracefulShutdown() {
	slog.Info("received interrupt")

	ctxTo, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// > program doesn't exit and waits instead for Shutdown to return.
	if err := s.httpSrv.Shutdown(ctxTo); err != nil {
		slog.Error("http shutdown", "error", err)
	}

	slog.Info("server terminated")
}

// serve the http server
func (s *Server) serve(ctx context.Context) error {
	handler, notifier, err := setupRoutes(ctx, &s.ready, s.JWTSecret)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", s.Port),
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
		slog.Error("listen and serve", "error", srv.ListenAndServe())
	}()

	s.httpSrv = srv
	return nil
}

func (s *Server) setReady() {
	s.ready.Store(true) // setup done
}

func setupRoutes(
	ctx context.Context,
	ready *atomic.Bool,
	secret string,
) (http.Handler, ws.Notifier, error) {
	mux := http.NewServeMux()

	notifier, err := ws.New(ctx, secret)
	if err != nil {
		return mux, notifier, err
	}

	mux.HandleFunc("GET /healthz", handlers.Health())
	mux.HandleFunc("GET /readyz", handlers.Ready(ready))
	mux.HandleFunc("GET /ws", notifier.Handler())
	mux.HandleFunc("GET /client", homeHandler())

	stack := middleware.Stack(
		middleware.Log,
	)

	return stack(mux), notifier, nil
}
