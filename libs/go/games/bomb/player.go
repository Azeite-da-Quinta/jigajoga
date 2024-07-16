package bomb

// GetState build the visible state by the player
func (g *Game) GetState(id int64) (PlayerState, error) {
	idx := g.getPlayerIdx(id)
	if idx == missing {
		return PlayerState{}, ErrPlayerNotFound
	}

	return PlayerState{
		Cards:         g.getPlayerCards(idx),
		Players:       g.ids,
		Readies:       g.readies,
		Playing:       g.playing,
		Round:         g.round,
		CardsRevealed: g.revealed,
		DefusesFound:  g.defusesFound,
		Role:          g.roles[idx],
		State:         g.state,
		Winner:        g.winner,
	}, nil
}

// TODO this state has to be pratical on the Front End.
// Stuff might need reorganizing or precalculating

// PlayerState contains the state visible by the player
type PlayerState struct {
	Cards         []Card    // shuffled hand of the player
	Players       []int64   // all players' ids
	Readies       []bool    // who's ready
	Playing       int64     // who's playing
	Round         uint8     // current round
	CardsRevealed uint8     // cards revealed this round / turn counter
	DefusesFound  uint8     // win condition of Heroes
	Role          Role      // team of the player
	State         StateKind // party's state
	Winner        Role
}
