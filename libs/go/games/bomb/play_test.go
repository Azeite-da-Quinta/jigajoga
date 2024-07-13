package bomb

import (
	"testing"
)

//revive:disable

var testPlayerIDs = []int64{1001, 1002, 1003, 1004}

func TestGame_Play(t *testing.T) {
	type fields struct {
		cards    []Card
		roles    []Role
		readies  []bool
		ids      []int64
		playing  int64
		round    uint8
		revealed uint8
		state    StateKind
		winner   Role
	}
	type args struct {
		id int64
		m  Move
	}
	type wants struct {
		err    bool
		state  StateKind
		winner Role
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   wants
	}{
		{
			name: "pick a card",
			fields: fields{
				state: Running,
				ids:   testPlayerIDs,
				cards: []Card{
					Bomb, Defuse, Defuse, Defuse, Defuse,
					Safe, Safe, Safe, Safe, Safe,
					Safe, Safe, Safe, Safe, Safe,
					Safe, Safe, Safe, Safe, Safe,
				},
				playing: testPlayerIDs[0],
			},
			args: args{
				id: testPlayerIDs[0],
				m: Move{
					Target:     testPlayerIDs[1],
					CardToDraw: 1,
				},
			},
			want: wants{
				err:   false,
				state: Running,
			},
		},
		{
			name: "pick the bomb",
			fields: fields{
				state: Running,
				ids:   testPlayerIDs,
				cards: []Card{
					Bomb, Defuse, Defuse, Defuse,
					Safe, Safe, Safe, Safe,
				},
				playing: testPlayerIDs[2],
			},
			args: args{
				id: testPlayerIDs[2],
				m: Move{
					Target:     testPlayerIDs[0],
					CardToDraw: 0,
				},
			},
			want: wants{
				err:    false,
				state:  Over,
				winner: Vilain,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				cards:    tt.fields.cards,
				roles:    tt.fields.roles,
				readies:  tt.fields.readies,
				ids:      tt.fields.ids,
				playing:  tt.fields.playing,
				round:    tt.fields.round,
				revealed: tt.fields.revealed,
				state:    tt.fields.state,
				winner:   tt.fields.winner,
			}

			err := g.Play(tt.args.id, tt.args.m)
			if (err != nil) != tt.want.err {
				t.Errorf("Game.Play() error = %v, wantErr %v",
					err, tt.want.err)
			}

			if g.state != tt.want.state {
				t.Errorf("Game.Play() state = %v, want %v",
					g.state, tt.want.state)
			}

			if g.winner != tt.want.winner {
				t.Errorf("Game.Play() winner = %v, want %v",
					g.winner, tt.want.winner)
			}
		})
	}
}

func TestGame_nextTurn(t *testing.T) {
	type fields struct {
		cards        []Card
		roles        []Role
		readies      []bool
		ids          []int64
		playing      int64
		round        uint8
		revealed     uint8
		defusesFound uint8
		state        StateKind
		winner       Role
	}
	type args struct {
		id int64
	}
	type wants struct {
		playing  int64
		revealed uint8
		round    uint8
		state    StateKind
		winner   Role
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   wants
	}{
		{
			name: "first play",
			fields: fields{
				ids:     testPlayerIDs,
				playing: testPlayerIDs[0],
			},
			args: args{
				id: testPlayerIDs[1],
			},
			want: wants{
				revealed: 1,
				playing:  testPlayerIDs[1],
			},
		},
		{
			name: "last card revelead first round",
			fields: fields{
				ids:      testPlayerIDs,
				playing:  testPlayerIDs[0],
				revealed: 3,
			},
			args: args{
				id: testPlayerIDs[1],
			},
			want: wants{
				revealed: 0,
				playing:  testPlayerIDs[1],
				round:    1,
			},
		},
		{
			name: "last card revelead last round",
			fields: fields{
				ids:      testPlayerIDs,
				playing:  testPlayerIDs[0],
				revealed: 3,
				round:    3,
			},
			args: args{
				id: testPlayerIDs[1],
			},
			want: wants{
				revealed: 0, // we don't care anymore
				playing:  testPlayerIDs[1],
				round:    4, // we don't care anymore
				winner:   Vilain,
				state:    Over,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				cards:        tt.fields.cards,
				roles:        tt.fields.roles,
				readies:      tt.fields.readies,
				ids:          tt.fields.ids,
				playing:      tt.fields.playing,
				round:        tt.fields.round,
				revealed:     tt.fields.revealed,
				defusesFound: tt.fields.defusesFound,
				state:        tt.fields.state,
				winner:       tt.fields.winner,
			}
			g.nextTurn(tt.args.id)

			if g.playing != tt.want.playing {
				t.Errorf("the hand should have passed %v want %v", g.playing, tt.want.playing)
			}

			if g.revealed != tt.want.revealed {
				t.Errorf("cards revealed not right %v want %v", g.revealed, tt.want.revealed)
			}

			if g.round != tt.want.round {
				t.Errorf("not the correct round %v want %v", g.round, tt.want.round)
			}

			if g.state != tt.want.state {
				t.Errorf("not the correct state %v want %v", g.state, tt.want.state)
			}

			if g.winner != tt.want.winner {
				t.Errorf("not the correct winner %v want %v", g.winner, tt.want.winner)
			}
		})
	}
}
