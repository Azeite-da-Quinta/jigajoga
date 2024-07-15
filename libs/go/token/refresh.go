package token

// Refresh token payload
type Refresh struct {
	// identifies the player
	PlayerIDField string `json:"player_id"`
	// increment-only value to avoid reusing refresh token
	Rotation int16 `json:"rotation"`
	// TODO: instead of this value, we could use JWT jti value
	// https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.7
}
