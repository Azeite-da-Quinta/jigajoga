package srv

import (
	"log/slog"
	"net/http"
)

func homeHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("http: home called")
		if r.Method != http.MethodGet {
			http.Error(
				w,
				http.StatusText(http.StatusMethodNotAllowed),
				http.StatusMethodNotAllowed,
			)
			return
		}

		http.ServeFile(w, r, "./dist/home.html")
	}
}
