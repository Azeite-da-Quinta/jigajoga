package token

// Access token payload
type Access struct {
	RoomIDField string `json:"r"`
	// identifies the player
	IDField   string `json:"id"`
	NameField string `json:"n"`
}
