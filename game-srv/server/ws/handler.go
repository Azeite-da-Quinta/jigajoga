package ws

import (
	"log/slog"
	"net/http"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/server/user"
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

// Handler creates the ws handler func
func (n *Notifier) Handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("ws: failed to ugprade", slog.String("err", err.Error()))
			return
		}

		// TODO: use real implementation
		t := user.MockToken()
		// client will close conn
		n.join(r.Context(), conn, t)
	}
}
