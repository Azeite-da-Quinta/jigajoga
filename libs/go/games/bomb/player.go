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

// Move is a play action
type Move struct {
	Target int64
}

// Play a move
func (g *Game) Play(id int64, m Move) error {
	if g.state == Lobby {
		return ErrPartyNotReady
	}
	// TODO check party over all well

	if g.playing != id {
		return ErrNotYourTurn
	}

	if g.getPlayerIdx(id) == missing ||
		g.getPlayerIdx(m.Target) == missing {
		return ErrPlayerNotFound
	}

	// TODO check if defuse

	// TODO check if bomb

	if areAllDefuseFound(g.playersCount(), 0) {
		g.state = Over
		g.winner = Hero
		return nil
	}

	g.nextTurn(m.Target)

	return nil
}

func (g *Game) nextTurn(id int64) {
	g.revealed++

	if isRoundEnd(g.playersCount(), int(g.revealed)) {
		g.nextRound()
	}

	g.playing = id //
}

func (g *Game) nextRound() {
	g.round++

	if areRoundsOver(g.round) {
		g.state = Over
		g.winner = Vilain
	}
}
