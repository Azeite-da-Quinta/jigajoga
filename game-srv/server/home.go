package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func homeHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("home called")
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		entries, err := os.ReadDir(".")
		if err != nil {
			slog.Error("read dir", "err", err)
		}

		for _, e := range entries {
			fmt.Println(e.Name())
		}

		http.ServeFile(w, r, "./dist/home.html")
	}
}
