package payload

// Message payload
type Message struct {
	Action  string `json:"a"`
	Content string `json:"c"`
	Target  string `json:"t"`
}
