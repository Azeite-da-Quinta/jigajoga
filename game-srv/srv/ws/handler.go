// Package ws is based on gorilla websockets.
package ws

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/party/event"
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/srv/ws/client"
	"github.com/Azeite-da-Quinta/jigajoga/libs/token"
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

// New returns a Notifier and runs the party Router
func New(n Notifier, secretKey string) (Handler, error) {
	b, err := token.Base64ToKey(secretKey)
	if err != nil {
		return Handler{}, err
	}

	return Handler{
		codec:    token.Codec{Key: b},
		Notifier: n,
	}, nil
}

// Handler creates the http handler that upgrades conn to ws
func (h *Handler) Handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := h.validateToken(r)
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

		cl := client.New(t)

		h.Notify(t, func(reply <-chan event.Reply) {
			cl.Pump(r.Context(), conn, reply)
		})
	}
}

// Handler builder
type Handler struct {
	Notifier
	codec token.Codec
}

// Notifier notifies the arrival of a new WS client
type Notifier interface {
	Notify(user.Token, func(reply <-chan event.Reply))
}

func (h *Handler) validateToken(r *http.Request) (user.Token, error) {
	const prefix = "Bearer "

	tok, err := h.codec.Decode(
		strings.TrimPrefix(
			r.Header.Get("Authorization"),
			prefix),
	)
	if err != nil {
		return nil, err
	}

	return user.FromToken(tok)
}
