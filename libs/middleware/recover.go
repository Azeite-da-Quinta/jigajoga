package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recover from a panic
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("recovered from panic", "value", r)
					fmt.Println(string(debug.Stack()))

					http.Error(w,
						http.StatusText(http.StatusServiceUnavailable),
						http.StatusServiceUnavailable,
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
}
