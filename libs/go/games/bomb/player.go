package bomb

// GetState build the visible state by the player
func (g *Game) GetState(id int64) (PlayerState, error) {
	idx := g.getPlayerIdx(id)
	if idx == missing {
		return PlayerState{}, ErrPlayerNotFound
	}

	return PlayerState{
		role:  g.roles[idx],
		cards: g.getPlayerCards(idx),
	}, nil
}

// PlayerState contains the state visible by the player
type PlayerState struct {
	cards []Card
	role  Role
}
