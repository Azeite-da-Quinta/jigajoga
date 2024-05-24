package ws

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// TODO
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Handler creates the ws handler func
func Handler(ctx context.Context, rt *Router) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("ws: failed to ugprade", slog.String("err", err.Error()))
			return
		}

		// TODO new func
		cl := client{id: "0", name: "n", inbox: make(chan []byte, 256)}
		// client will close conn
		cl.run(ctx, rt, conn)
	}
}
