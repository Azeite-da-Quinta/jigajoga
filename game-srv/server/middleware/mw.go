package middleware

import "net/http"

type Global func(http.Handler) http.Handler

func Stack(mws ...Global) Global {
	return func(next http.Handler) http.Handler {
		for i := range mws {
			f := mws[len(mws)-1-i]
			next = f(next)
		}

		return next
	}
}
