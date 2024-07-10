// Package kvstore around Redis
// References:
// https://redis.io/docs/latest/develop/interact/pubsub/
package kvstore

import (
	"context"
	"log/slog"
	"strings"

	"github.com/Azeite-da-Quinta/jigajoga/libs/go/envelope"
	"github.com/Azeite-da-Quinta/jigajoga/libs/go/slogt"
	"github.com/Azeite-da-Quinta/jigajoga/libs/go/user"
	"github.com/Azeite-da-Quinta/jigajoga/svc/game-srv/pkg/party/event"
	"github.com/Azeite-da-Quinta/jigajoga/svc/game-srv/pkg/remote"
	"github.com/redis/go-redis/v9"
)

// Client wrapper around Redis
type Client struct {
	rdb *redis.Client
}

// New initializes a new client
func New() Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return Client{
		rdb: rdb,
	}
}

const subscribeBuffer = 1024

const (
	roomsTopic = "rooms:*:*"
	// TODO replace messagesTopic by moves
	messagesTopic = "messages:*:*"
)

// Subscribe
/* func (c *Client) Subscribe(ctx context.Context) <-chan event.Envelope {
	out := make(chan event.Envelope, subscribeBuffer)

	pubsub := c.rdb.PSubscribe(ctx, rooms)
	ch := pubsub.Channel()
	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				slog.Info("redis: rooms subscriber closed")
				return
			case msg := <-ch:
				if envelope.IsTopic(msg.Channel) {
					e, err := envelope.FromBytes([]byte(msg.Payload))
					if err != nil {
						slog.Error("redis: failed to parse bytes",
							"channel", msg.Channel)
						continue
					}
					fmt.Println("read message from Redis")

					out <- e
					return
				} else if notifier.IsTopic(msg.Channel) {
					fmt.Println("read notification from Redis")

					fmt.Println("msg", msg.Channel, msg.Payload, msg.PayloadSlice)
					e, err := notifier.FromBytes([]byte(msg.Payload))
					if err != nil {
						slog.Error("redis: failed to parse bytes",
							"channel", msg.Channel)
						continue
					}
					fmt.Println(e)
				}

				slog.Error("redis: unknown message format",
					"channel", msg.Channel)
			}
		}
	}()

	return out
}
*/

// SubRooms subcribes to events
func (c *Client) SubRooms(ctx context.Context) <-chan remote.Payload {
	out := make(chan remote.Payload, subscribeBuffer)

	pubsub := c.rdb.PSubscribe(ctx, roomsTopic)
	go roomsWorker(ctx, pubsub.Channel(), out)

	return out
}

func roomsWorker(
	ctx context.Context,
	ch <-chan *redis.Message,
	out chan<- remote.Payload,
) {
	defer close(out)

	for {
		select {
		case <-ctx.Done():
			slog.Info("redis: rooms subscriber closed")
			return
		case msg := <-ch:
			e := getEvent(msg.Channel)

			t, err := user.Deserialize([]byte(msg.Payload))
			if err != nil {
				slog.Warn("pubsub: deserialize", slogt.Error(err))
			}

			out <- Payload{event: e, Token: t}
		}
	}
}

// SubMessages subcribes to message events
func (c *Client) SubMessages(ctx context.Context) <-chan envelope.Message {
	out := make(chan envelope.Message, subscribeBuffer)

	pubsub := c.rdb.PSubscribe(ctx, messagesTopic)
	go func(
		ctx context.Context,
		ch <-chan *redis.Message,
		out chan<- envelope.Message,
	) {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				slog.Info("redis: messages subscriber closed")
				return
			case msg := <-ch:
				m, err := envelope.FromBytes([]byte(msg.Payload))
				if err != nil {
					slog.Warn("pubsub: deserialize", slogt.Error(err))
				}

				out <- m
			}
		}
	}(ctx, pubsub.Channel(), out)

	return out
}

// Publish an event.Query
func (c *Client) Publish(ctx context.Context, e event.Envelope) {
	b, err := e.Serialize()
	if err != nil {
		slog.Error("redis: failed to serialize envelope", slogt.Error(err))
	}

	err = c.rdb.Publish(ctx, e.Topic(), b).Err()
	if err != nil {
		slog.Error("redis: publisher", slogt.Error(err))
	}
}

// getEvent extracts the event from the topic
func getEvent(s string) string {
	const formatSize = 3

	array := strings.Split(s, ":")
	if len(array) != formatSize {
		return ""
	}

	const last = 2
	return array[last]
}
