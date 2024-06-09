// Package ws is based on gorilla websockets.
package ws

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/Azeite-da-Quinta/jigajoga/libs/user"
	"github.com/gorilla/websocket"
)

const (
	readBufSize  = 1024
	writeBufSize = readBufSize
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  readBufSize,
	WriteBufferSize: writeBufSize,
	// TODO
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Handler creates the http handler that upgrades conn to ws
func (n *Notifier) Handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := n.validateToken(r)
		if err != nil {
			// we could distinguish between Forbidden
			// vs Unauthorized 🤷
			http.Error(w,
				http.StatusText(http.StatusForbidden),
				http.StatusForbidden,
			)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("ws: failed to ugprade", slog.String("err", err.Error()))
			return
		}

		// client will close conn
		n.join(r.Context(), conn, t)
	}
}

func (n *Notifier) validateToken(r *http.Request) (user.Token, error) {
	const prefix = "Bearer "

	token, err := n.codec.Decode(
		strings.TrimPrefix(
			r.Header.Get("Authorization"),
			prefix),
	)
	if err != nil {
		return nil, err
	}

	return user.FromToken(token)
}
