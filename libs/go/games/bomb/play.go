package bomb

// Move is a play action
type Move struct {
	Target     int64
	CardToDraw uint8
}

// Play a move
func (g *Game) Play(id int64, m Move) error {
	pIdx, err := g.checkPreconditions(id, m)
	if err != nil {
		return err
	}

	c, err := g.drawCard(pIdx, int(m.CardToDraw))
	if err != nil {
		return err
	}

	// check victory conditions
	switch c {
	case Bomb:
		g.state = Over
		g.winner = Vilain
		return nil
	case Defuse:
		g.defusesFound++

		if areAllDefuseFound(g.playersCount(), int(g.defusesFound)) {
			g.state = Over
			g.winner = Hero
			return nil
		}
	}

	isNextRound := g.nextTurn(m.Target)
	if isNextRound {
		g.setNextRoundCards()
	}

	return nil
}

func (g Game) checkPreconditions(id int64, m Move) (int, error) {
	switch g.state {
	case Lobby:
		return 0, ErrPartyNotReady
	case Over:
		return 0, ErrPartyOver
	}

	if g.playing != id {
		return 0, ErrNotYourTurn
	}

	if m.Target == g.playing {
		return 0, ErrCannotDrawSelf
	}

	pIdx := g.getPlayerIdx(m.Target)

	if g.getPlayerIdx(id) == missing ||
		pIdx == missing {
		return 0, ErrPlayerNotFound
	}

	return pIdx, nil
}

func (g *Game) nextTurn(id int64) bool {
	g.revealed++
	g.playing = id

	if isRoundEnd(g.playersCount(), int(g.revealed)) {
		g.nextRound()
		return true
	}

	return false
}

func (g *Game) nextRound() {
	g.revealed = 0
	g.round++

	if areRoundsOver(g.round) {
		g.state = Over
		g.winner = Vilain
	}
}
