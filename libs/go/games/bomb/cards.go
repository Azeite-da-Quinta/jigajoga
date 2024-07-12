package bomb

import (
	"math/rand"
)

// Card type
type Card uint8

// enum kinds of cards
const (
	Unkown Card = iota
	Safe
	Defuse
	Bomb
)

func (g *Game) drawCard(idx int) {}

// getPlayerCards returns the cards a player has
func (g Game) getPlayerCards(idx int) []Card {
	num := cardsPerRound(int(g.round))
	begin := idx * num

	cards := g.cards[begin : begin+num]

	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	return cards
}

// genCards generates the initial deck of cards and
// shuffles it
func genCards(count int) []Card {
	safe, defuse, bomb := initialCardsCount(count)

	cards := make([]Card, 0, safe+defuse+bomb)

	for range safe {
		cards = append(cards, Safe)
	}

	for range defuse {
		cards = append(cards, Defuse)
	}

	cards = append(cards, Bomb)

	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	return cards
}

// initialCardsCount returns how many of each card type
// the deck should start with.
func initialCardsCount(count int) (safe, defuse, bomb int) {
	bomb = 1

	switch count {
	case 4:
		safe = 15
		defuse = 4
	case 5:
		safe = 19
		defuse = 5
	case 6:
		safe = 23
		defuse = 6
	case 7:
		safe = 27
		defuse = 7
	case 8:
		safe = 31
		defuse = 8
	}

	return safe, defuse, bomb
}
