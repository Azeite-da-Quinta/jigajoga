package notifier

import (
	"fmt"
	"strings"

	"github.com/Azeite-da-Quinta/jigajoga/libs/go/user"
	"github.com/Azeite-da-Quinta/jigajoga/svc/game-srv/pkg/party/event"
)

const topic = "rooms"

// IsTopic
func IsTopic(s string) bool {
	if len(s) < len(topic) {
		return false
	}

	return strings.Contains(s, topic)
}

type join struct {
	user.Token
	// Reply in this channel with the room's write channel
	reply chan event.Reply
}

func (j join) Get() user.Token {
	return j.Token
}

func (j join) Topic() string {
	return prepareTopic(int64(j.RoomID()), "join")
}

func (j join) Reply() chan event.Reply {
	return j.reply
}

type leave struct {
	user.Token
}

func (l leave) Get() user.Token {
	return l.Token
}

func (l leave) Topic() string {
	return prepareTopic(int64(l.RoomID()), "leave")
}

// TODO this function could be implemented elsewhere
func prepareTopic(roomID int64, e string) string {
	return fmt.Sprintf("%s:%d:%s", topic, roomID, e)
}
