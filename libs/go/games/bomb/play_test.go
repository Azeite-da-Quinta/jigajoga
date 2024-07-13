package bomb

import "testing"

//revive:disable
func TestGame_Play(t *testing.T) {
	players := []int64{11, 12, 13, 14}

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
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantState StateKind
	}{
		{
			name: "pick a card",
			fields: fields{
				state:   Running,
				ids:     players,
				cards:   []Card{},
				playing: players[0],
			},
			args: args{
				id: players[0],
				m: Move{
					Target: players[1],
				},
			},
			wantErr:   false,
			wantState: Running,
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
			if (err != nil) != tt.wantErr {
				t.Errorf("Game.Play() error = %v, wantErr %v",
					err, tt.wantErr)
			}

			if g.state != tt.wantState {
				t.Errorf("Game.Play() state = %v, want %v",
					g.state, tt.wantState)
			}
		})
	}
}
