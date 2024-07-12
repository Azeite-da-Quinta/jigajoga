package bomb

import "math/rand"

// Role type
type Role uint8

// enum kinds of players
const (
	None Role = iota
	Hero
	Vilain
)

func genRoles(count int) []Role {
	heroes, vilains := getRolesCount(count)

	roles := make([]Role, 0, heroes+vilains)

	for range heroes {
		roles = append(roles, Hero)
	}

	for range vilains {
		roles = append(roles, Vilain)
	}

	rand.Shuffle(len(roles), func(i, j int) {
		roles[i], roles[j] = roles[j], roles[i]
	})

	// with 4 and 7 players, we discard one card
	if len(roles) > count {
		roles = roles[:count]
	}

	return roles
}

func getRolesCount(count int) (heroes, vilains int) {
	switch count {
	case 4, 5:
		heroes = 3
		vilains = 2
	case 6:
		heroes = 4
		vilains = 2
	case 7, 8:
		heroes = 5
		vilains = 3
	}

	return heroes, vilains
}
