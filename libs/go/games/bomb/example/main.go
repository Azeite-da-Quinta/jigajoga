package main

import (
	"fmt"

	"github.com/Azeite-da-Quinta/jigajoga/libs/go/games/bomb"
)

// go run libs/go/games/bomb/example/main.go
func main() {
	players := []int64{
		100, 101, 102, 103, 104, 105, 106, 107,
	}

	game := bomb.New()

	// All players join the game
	for _, id := range players {
		game.Join(id)
	}

	// All players mark as ready
	for _, id := range players {
		game.SetReady(id, true)
	}

	// obtain the game state of a player
	state, err := game.GetState(players[0])
	if err != nil {
		panic(err)
	}

	// print the cards of the player
	fmt.Println(state.Cards)

	err = game.Play(state.Playing, bomb.Move{
		Target:     players[7],
		CardToDraw: 2,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("played a move")
}
