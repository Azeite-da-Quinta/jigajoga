package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type StatusWriter struct {
	http.ResponseWriter
	Status int
}

// WriteHeader wraps underlying WriterHeader and keeps http status code
func (sw *StatusWriter) WriteHeader(status int) {
	sw.ResponseWriter.WriteHeader(status)
	sw.Status = status
}

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
