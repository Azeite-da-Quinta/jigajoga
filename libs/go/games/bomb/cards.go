package bomb

import (
	"math/rand"
	"slices"
)

// Card type
type Card uint8

// enum kinds of cards
const (
	Unkown Card = iota
	// No effect card
	Safe
	// Victory condition
	Defuse
	// Game over condition
	Bomb
	// Card already drawn another round
	Drawn
)

func (g *Game) drawCard(playerIdx, cardIdx int) (Card, error) {
	idx, err := g.cardIdx(playerIdx, cardIdx)
	if err != nil {
		return Unkown, err
	}

	c := g.cards[idx]
	switch c {
	case Unkown, Drawn:
		return c, ErrCardNotFound
	}

	// We flag the card so it's not drawn
	// next turns
	g.cards[idx] = Drawn

	return c, nil
}

func (g Game) cardIdx(playerIdx, cardIdx int) (idx int, err error) {
	num := cardsPerRound(int(g.round))

	if cardIdx < 0 ||
		cardIdx >= num {
		return 0, ErrCardNotFound
	}

	begin := playerIdx * num
	if begin+cardIdx >= len(g.cards) {
		return 0, ErrCardNotFound
	}

	return begin + cardIdx, nil
}

// getPlayerCards returns the cards a player has.
// internally called, so shouldn't OOB
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
	//revive:disable:add-constant
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

func (g *Game) setNextRoundCards() {
	revealedCards := g.playersCount()

	slices.Sort(g.cards)

	g.cards = g.cards[:len(g.cards)-revealedCards]

	rand.Shuffle(len(g.cards), func(i, j int) {
		g.cards[i], g.cards[j] = g.cards[j], g.cards[i]
	})
}
