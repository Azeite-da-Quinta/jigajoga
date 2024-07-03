package srv

import (
	"context"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/kvstore"
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/party"
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/party/notifier"
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/remote"
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/srv/ws"
)

func (s *Server) setupWS(ctx context.Context) (ws.Handler, error) {
	n, routed, spy := notifier.NewTee(ctx)

	h, err := ws.New(n, s.JWTSecret)
	if err != nil {
		return h, err
	}
	s.closers = append(s.closers, n)

	cl := kvstore.New()

	multiplexer := remote.New(n, spy)
	go multiplexer.Run(ctx, &cl)

	// go cl.Subscribe(ctx, n)
	// go cl.Publisher(ctx, right)

	rt := party.NewRouter(routed)
	go rt.Run(ctx)

	return h, nil
}
