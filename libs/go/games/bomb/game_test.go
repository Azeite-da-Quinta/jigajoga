package bomb

import (
	"testing"
)

//revive:disable:add-constant
func Test(t *testing.T) {
	players := []int64{11, 12, 13, 14}

	g := New()

	for _, id := range players {
		g.Join(id)
	}

	if len(g.ids) != len(players) {
		t.Error("players missing")
		return
	}

	for i := range 2 {
		g.SetReady(players[i], true)
	}

	if g.state != Lobby {
		t.Error("players should not be ready")
		return
	}

	for i := 2; i < 4; i++ {
		g.SetReady(players[i], true)
	}

	if g.state != Running {
		t.Error("players should all be ready")
	}
}
