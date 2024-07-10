// Package main contains examples for Redis
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/redis/go-redis/v9"
)

const channel = "abcdef"

//revive:disable
func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		pubsub := rdb.Subscribe(ctx, channel)
		defer pubsub.Close()

		ch := pubsub.Channel()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("sub done")
				return
			case msg := <-ch:
				fmt.Println(msg.Channel, msg.Payload)
			}
		}
	}()

	go func() {
		for range 10 {
			err := rdb.Publish(ctx, channel, "payload").Err()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	<-ctx.Done()
	fmt.Println("bye")
	time.Sleep(10 * time.Millisecond)
}
