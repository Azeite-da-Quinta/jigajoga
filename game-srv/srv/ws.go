package srv

import (
	"context"
	"log/slog"

	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/kvstore"
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/party"
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/party/notifier"
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/pkg/remote"
	"github.com/Azeite-da-Quinta/jigajoga/game-srv/srv/ws"
)

func (s *Server) setupWS(ctx context.Context) (ws.Handler, error) {
	switch s.Mode {
	case string(single):
		return s.setupSingle(ctx)
	case string(cluster):
		return s.setupCluster(ctx)
	default:
		slog.Warn("srv: setup mode not implemented")
	}

	return s.setupSingle(ctx)
}

func (s *Server) setupSingle(ctx context.Context) (ws.Handler, error) {
	n, routed := notifier.New(ctx)

	h, err := ws.New(n, s.JWTSecret)
	if err != nil {
		return h, err
	}
	s.closers = append(s.closers, n)

	rt := party.NewRouter(routed)
	go rt.Run(ctx)

	return h, nil
}

func (s *Server) setupCluster(ctx context.Context) (ws.Handler, error) {
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
