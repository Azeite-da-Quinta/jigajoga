package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/libs/go/slogt"
)

func urlReady(conf Config) string {
	return fmt.Sprintf("http://%s/readyz", conf.Host)
}

func httpReady(ctx context.Context, conf Config) error {
	c := http.Client{Timeout: time.Second}

	resp, err := c.Get(urlReady(conf))
	if err != nil {
		slog.Error("failed to http get ready", slogt.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.LogAttrs(ctx, slog.LevelError, "status not ok",
			slog.Int("status", resp.StatusCode))
		return errors.New("server not ready")
	}

	return nil
}
