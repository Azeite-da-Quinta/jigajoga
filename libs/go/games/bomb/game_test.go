package bomb

import (
	"testing"
)

//revive:disable:add-constant
func Test(t *testing.T) {
	g := New()

	for _, id := range testPlayerIDs {
		g.Join(id)
	}

	if len(g.ids) != len(testPlayerIDs) {
		t.Error("players missing")
		return
	}

	for i := range 2 {
		g.SetReady(testPlayerIDs[i], true)
	}

	if g.state != Lobby {
		t.Error("players should not be ready")
		return
	}

	for i := 2; i < 4; i++ {
		g.SetReady(testPlayerIDs[i], true)
	}

	if g.state != Running {
		t.Error("players should all be ready")
	}
}
