// Package handlers has functions that return http.Handler for misc stuff
package handlers

import (
	"net/http"
)

// Health dummy healthy handler
func Health() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy"))
	}
}
