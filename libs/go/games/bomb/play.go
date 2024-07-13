package bomb

// Move is a play action
type Move struct {
	Target     int64
	CardToDraw uint8
}

// TODO check that player doesn't draw from himself

// Play a move
func (g *Game) Play(id int64, m Move) error {
	if g.state == Lobby {
		return ErrPartyNotReady
	}
	// TODO check party over all well

	if g.playing != id {
		return ErrNotYourTurn
	}

	t := g.getPlayerIdx(m.Target)

	if g.getPlayerIdx(id) == missing ||
		t == missing {
		return ErrPlayerNotFound
	}

	c, err := g.drawCard(t, int(m.CardToDraw))
	if err != nil {
		return err
	}

	if c == Bomb {
		g.state = Over
		g.winner = Vilain
		return nil
	}

	if c == Defuse {
		g.defusesFound++
	}

	if areAllDefuseFound(g.playersCount(), int(g.defusesFound)) {
		g.state = Over
		g.winner = Hero
		return nil
	}

	isNextRound := g.nextTurn(m.Target)
	if isNextRound {
		g.setNextRoundCards()
	}

	return nil
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
