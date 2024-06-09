// Package client to test the server
package client

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/libs/token"
	"github.com/Azeite-da-Quinta/jigajoga/libs/user"
)

// Config of Dial
type Config struct {
	Version   string
	Host      string
	JWTSecret string
	NbWorkers int
	NbWrites  int
}

// Dial connects to the server with N workers doing N jobs (write messages)
func Dial(conf Config) {
	ctxS, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ctx, cancel := context.WithTimeout(ctxS, 1*time.Minute)
	defer cancel()

	if err := httpReady(ctx, conf); err != nil {
		return
	}

	doJobs(ctx, conf)

	slog.Info("closing clients")
	time.Sleep(100 * time.Millisecond)
}

// dev
const node = 0

func doJobs(ctx context.Context, conf Config) {
	var wg sync.WaitGroup
	wg.Add(conf.NbWorkers)

	url := urlWS(conf.Host)

	b, err := token.Base64ToKey(conf.JWTSecret)
	if err != nil {
		panic(err)
	}

	cod := token.Codec{
		Key: b,
	}

	fac, err := user.NewFactory(node)
	if err != nil {
		panic(err)
	}
	room := fac.NewRoom()

	for i := range conf.NbWorkers {
		claims := fac.NewUser(mockName(i), room).
			ToToken().
			Claims(time.Now())

		jwt, err := cod.Encode(claims)
		if err != nil {
			panic(err)
		}

		w := worker{
			url:      url,
			jwt:      jwt,
			nbWrites: conf.NbWrites,
			num:      i,
		}

		go w.run(ctx, &wg)
	}

	wg.Wait()
}

func urlWS(host string) string {
	return fmt.Sprintf("ws://%s/ws", host)
}

func mockName(i int) string {
	return fmt.Sprintf("bob-the-builder-%d", i)
}
