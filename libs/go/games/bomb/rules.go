package bomb

const (
	minPlayersToStart = 4
	maxPlayers        = 8
	maxRounds         = 4 // vilains winning condition, all 4 rounds played
	startingCards     = 5 // how many cards at the start of the game
)

func hasMinplayers(count int) error {
	if count < minPlayersToStart {
		return ErrMinPlayers
	}

	return nil
}

func hasMaxPlayers(count int) error {
	if count == maxPlayers {
		return ErrTooManyPlayers
	}

	return nil
}

func cardsPerRound(round int) int {
	// each round, players start with one card less
	return startingCards - round
}

func isRoundEnd(playersCount, cardsDealt int) bool {
	return playersCount >= cardsDealt
}

func areRoundsOver(round uint8) bool {
	return round > maxRounds
}

func areAllDefuseFound(playersCount, defuse int) bool {
	return playersCount == defuse
}
