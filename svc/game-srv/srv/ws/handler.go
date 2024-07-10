// Package ws is based on gorilla websockets.
package ws

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Azeite-da-Quinta/jigajoga/libs/go/envelope"
	"github.com/Azeite-da-Quinta/jigajoga/libs/go/token"
	"github.com/Azeite-da-Quinta/jigajoga/libs/go/user"
	"github.com/Azeite-da-Quinta/jigajoga/svc/game-srv/pkg/party/event"
	"github.com/Azeite-da-Quinta/jigajoga/svc/game-srv/srv/ws/client"
	"github.com/gorilla/websocket"
)

const (
	readBufSize  = 1024
	writeBufSize = readBufSize
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  readBufSize,
	WriteBufferSize: writeBufSize,
	Subprotocols:    []string{envelope.Subprotocol},
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
		token, err := getToken(r)
		if err != nil {
			slog.Info("ws: failed to retrieve token", slog.String("err", err.Error()))

			http.Error(w,
				http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized,
			)
			return
		}

		t, err := h.validateToken(token)
		if err != nil {
			slog.Info("ws: rejected token", slog.String("err", err.Error()))

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

func getToken(r *http.Request) (string, error) {
	const prefix = "base64url.bearer.authorization."

	for _, sub := range websocket.Subprotocols(r) {
		if strings.HasPrefix(sub, prefix) {
			return strings.TrimPrefix(sub, prefix), nil
		}
	}

	return "", errors.New("bearer token not found")
}

func (h *Handler) validateToken(t string) (user.Token, error) {
	tok, err := h.codec.Decode(t)
	if err != nil {
		return nil, err
	}

	return user.FromToken(tok)
}
