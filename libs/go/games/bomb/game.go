package bomb

import (
	"slices"
)

// New creates a Bomb game
func New() Game {
	return Game{
		ids:     make([]int64, 0, maxPlayers),
		readies: make([]bool, 0, maxPlayers),
		roles:   make([]Role, 0, maxPlayers),
	}
}

// Game contains the state
type Game struct {
	cards        []Card
	roles        []Role
	readies      []bool
	ids          []int64
	playing      int64 // id of the player currently playing
	round        uint8 // current round
	revealed     uint8 // cards revealed this round / turn counter
	defusesFound uint8 // win condition of Heroes
	state        StateKind
	winner       Role
}

const missing = -1

// Join adds the player to the game
func (g *Game) Join(id int64) error {
	if g.state != Lobby {
		return ErrPartyAlreadyStarted
	}

	if err := hasMaxPlayers(g.playersCount()); err != nil {
		return err
	}

	if g.getPlayerIdx(id) == missing {
		g.ids = append(g.ids, id)
		g.readies = append(g.readies, false)
	}

	return nil
}

// SetReady changes the status of readiness
func (g *Game) SetReady(id int64, b bool) {
	idx := g.getPlayerIdx(id)
	if idx == missing {
		return
	}

	g.readies[idx] = b
	g.checkReady()
}

func (g Game) getPlayerIdx(id int64) int {
	return slices.Index(g.ids, id)
}

func (g *Game) checkReady() {
	if g.state != Lobby {
		return
	}

	if hasMinplayers(g.playersCount()) != nil {
		return // Error is silenced.
	}

	waiting := slices.Contains(g.readies, false)
	if waiting {
		return
	}

	g.start()
}

func (g *Game) start() {
	g.state = Running

	count := g.playersCount()

	g.roles = genRoles(count)
	g.cards = genCards(count)

	// TODO define random starting player
}

func (g Game) playersCount() int {
	return len(g.ids)
}
