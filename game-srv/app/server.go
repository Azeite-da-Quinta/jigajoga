package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
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

	var ready atomic.Bool
	s := Serve(ctx, c.Port, &ready)

	// setup done
	ready.Store(true)

	// waits for app to interrupt
	<-ctx.Done()
	slog.Info("closing: received interrupt")

	ctxTo, cancel := context.WithTimeoutCause(ctx, 2*time.Second, errors.New("closing: timeout exceed"))
	defer cancel()

	s.Shutdown(ctxTo) // > program doesn't exit and waits instead for Shutdown to return.
	slog.Info("closing: server terminated")
}

func Serve(ctx context.Context, port int, ready *atomic.Bool) *http.Server {
	mux := http.NewServeMux()

	// TODO: move handlers elsewhere
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy"))
	})

	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		if ready == nil || !ready.Load() {
			http.Error(w,
				http.StatusText(http.StatusServiceUnavailable),
				http.StatusServiceUnavailable,
			)
		}
		w.Write([]byte("Realthy"))
	})

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		// uses the default slog as log
		ErrorLog:    slog.NewLogLogger(slog.Default().Handler(), slog.LevelInfo),
		BaseContext: func(l net.Listener) context.Context { return ctx },
	}

	// 🚀
	go slog.Error("listen and serve", "error", s.ListenAndServe())

	return s
}
