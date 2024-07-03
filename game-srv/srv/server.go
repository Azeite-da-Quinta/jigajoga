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
	Mode      string
	Port      int
}

type closer interface {
	Close()
}

// Server wraps an http server and more. It is setup by Config
type Server struct {
	httpSrv *http.Server
	Config
	closers []closer
	ready   atomic.Bool
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
		slog.Error("server: failed to serve http", slogt.Error(err))
		panic(err)
	}
	s.ready.Store(true) // setup done

	// waits for app to interrupt
	<-ctx.Done()
	s.gracefulShutdown()
}

func (s *Server) gracefulShutdown() {
	slog.Info("server: received interrupt")

	s.ready.Store(false) // server is no longer ready

	ctxTo, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// > program doesn't exit and waits instead for Shutdown to return.
	if err := s.httpSrv.Shutdown(ctxTo); err != nil {
		slog.Error("server: http shutdown", "error", err)
	}

	slog.Info("server: terminated")
}

// serve the http server
func (s *Server) serve(ctx context.Context) error {
	handler, err := s.setupRoutes(ctx, &s.ready)
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
		defer func() {
			for _, c := range s.closers {
				c.Close()
			}
		}()
		slog.Error("server: listen and serve", "error", srv.ListenAndServe())
	}()

	s.httpSrv = srv
	return nil
}

func (s *Server) setupRoutes(
	ctx context.Context,
	ready *atomic.Bool,
) (http.Handler, error) {
	mux := http.NewServeMux()

	h, err := s.setupWS(ctx)
	if err != nil {
		return mux, err
	}

	mux.HandleFunc("GET /healthz", handlers.Health())
	mux.HandleFunc("GET /readyz", handlers.Ready(ready))
	mux.HandleFunc("GET /ws", h.Handler())
	mux.HandleFunc("GET /client", homeHandler())

	stack := middleware.Stack(
		middleware.Log,
	)

	return stack(mux), nil
}
