package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Stack(mws ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := range mws {
			f := mws[len(mws)-1-i]
			next = f(next)
		}

		return next
	}
}
