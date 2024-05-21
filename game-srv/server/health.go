package server

import (
	"net/http"
	"sync/atomic"
)

func healthHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy"))
	}
}

func readyHandler(ready *atomic.Bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if ready == nil || !ready.Load() {
			http.Error(w,
				http.StatusText(http.StatusServiceUnavailable),
				http.StatusServiceUnavailable,
			)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))
	}
}
