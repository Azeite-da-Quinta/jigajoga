package middleware

import (
	"bufio"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// StatusWriter wrapper around ResponseWriter that
// keeps track of the HTTP status code and
// implements http.Hijaker
type StatusWriter struct {
	http.ResponseWriter
	Status int
}

// WriteHeader wraps underlying WriterHeader and keeps http status code
func (sw *StatusWriter) WriteHeader(status int) {
	sw.ResponseWriter.WriteHeader(status)
	sw.Status = status
}

// Hijack implements http.Hijacker
func (sw *StatusWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := sw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

// Log http Handler captures elapsed time and logs to default slog
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			sw := StatusWriter{
				ResponseWriter: w,
				Status:         http.StatusOK, // default value
			}

			next.ServeHTTP(&sw, r)

			slog.Info("http req",
				slog.Int("status", sw.Status),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Duration("elapsed", time.Since(now)),
			)
		})
}
