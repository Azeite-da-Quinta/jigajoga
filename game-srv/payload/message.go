// Package payload contains the payloads to be exchanged in websockets
package payload

// Message payload
type Message struct {
	Action  string `json:"a"`
	Content string `json:"c"`
	Target  string `json:"t"`
}
