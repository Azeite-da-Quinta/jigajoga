package bomb

import (
	"errors"
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

	for i := range testPlayerIDs {
		g.SetReady(testPlayerIDs[i], true)
	}

	if g.state != Running {
		t.Error("players should all be ready")
	}

	if err := g.Join(1005); !errors.Is(err, ErrPartyAlreadyStarted) {
		t.Errorf("player joined after party was started or err wrong: %v", err)
	}

	moves := []Move{
		// round 0
		{Target: testPlayerIDs[1], CardToDraw: 0},
		{Target: testPlayerIDs[2], CardToDraw: 0},
		{Target: testPlayerIDs[3], CardToDraw: 0},
		{Target: testPlayerIDs[0], CardToDraw: 0},
		// round 1
		{Target: testPlayerIDs[1], CardToDraw: 0},
		{Target: testPlayerIDs[2], CardToDraw: 0},
		{Target: testPlayerIDs[3], CardToDraw: 0},
		{Target: testPlayerIDs[0], CardToDraw: 0},
		// round 2
		{Target: testPlayerIDs[1], CardToDraw: 0},
		{Target: testPlayerIDs[2], CardToDraw: 0},
		{Target: testPlayerIDs[3], CardToDraw: 0},
		{Target: testPlayerIDs[0], CardToDraw: 0},
		// round 3
		{Target: testPlayerIDs[1], CardToDraw: 0},
		{Target: testPlayerIDs[2], CardToDraw: 0},
		{Target: testPlayerIDs[3], CardToDraw: 0},
		{Target: testPlayerIDs[0], CardToDraw: 0},
	}
	g.playing = testPlayerIDs[0] // force for test
	for _, m := range moves {
		// Code to skip bomb (:
		pIdx := g.getPlayerIdx(g.playing)
		cIdx, err := g.cardIdx(pIdx, int(m.CardToDraw))
		if err != nil {
			t.Errorf("card idx: %v", err)
		}
		c := g.cards[cIdx]
		if c == Bomb {
			m.CardToDraw++
		}

		err = g.Play(g.playing, m)
		if err != nil {
			t.Errorf("play failed: %v", err)
			return
		}

		if g.state == Over {
			return
		}
	}
	if g.winner == None {
		t.Errorf("no one won %v", g.winner)
	}
}
