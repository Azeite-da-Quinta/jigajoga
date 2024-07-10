// Package envelope contains the payloads to be exchanged in websockets
package envelope

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Subprotocol is the ws subprotocol to use
const Subprotocol = "v0.jigajoga.json"

const topic = "messages"

// IsTopic
func IsTopic(s string) bool {
	if len(s) < len(topic) {
		return false
	}

	return strings.Contains(s, topic)
}

// Message payload sent by a user
type Message struct {
	// Action  string `json:"a"`
	Content string `json:"c"`
	From    string `json:"f,omitempty"`
	To      string `json:"t"`
}

// Serialize the content to json bytes
func (m Message) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

// Topic returns the approriate topic for this payload
func (m Message) Topic() string {
	// TODO
	return fmt.Sprintf("%s:TODO:%s", topic, m.From)
}

// FromBytes a message in JSON bytes
func FromBytes(data []byte) (Message, error) {
	var m Message
	err := json.Unmarshal(data, &m)

	return m, err
}
