// Package middleware contains http.Handler middlewares
package middleware

import "net/http"

// Global defines a Global middleware to apply to a mux
type Global func(http.Handler) http.Handler

// Stack returns the stacked list of middlewares
func Stack(mws ...Global) Global {
	return func(next http.Handler) http.Handler {
		for i := range mws {
			// iterates from the back
			//revive:disable:add-constant
			f := mws[len(mws)-1-i]
			next = f(next)
		}

		return next
	}
}
