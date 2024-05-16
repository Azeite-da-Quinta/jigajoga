package client

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Config struct {
	Version string
	Host    string
}

func Dial(conf Config) {
	c := http.Client{Timeout: time.Duration(1) * time.Second}

	resp, err := c.Get(conf.Host + "/readyz")
	if err != nil {
		slog.Error("failed to http get ready", "error", err)
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to check ready", "error", err)
		return
	}
	fmt.Printf("response %s:\n%s", resp.Status, b)
}
