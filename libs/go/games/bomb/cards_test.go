package bomb

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

//revive:disable
func Test_genCards(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name string
		want []Card
		args args
	}{
		{
			name: "4 players",
			args: args{
				count: 4,
			},
			want: []Card{
				Bomb,
				Defuse, Defuse, Defuse, Defuse,
				Safe, Safe, Safe, Safe, Safe,
				Safe, Safe, Safe, Safe, Safe,
				Safe, Safe, Safe, Safe, Safe,
			},
		},
		{
			name: "8 players",
			args: args{
				count: 8,
			},
			want: []Card{
				Bomb,
				Defuse, Defuse, Defuse, Defuse,
				Defuse, Defuse, Defuse, Defuse,
				Safe, Safe, Safe, Safe, Safe, Safe, Safe, Safe, Safe, Safe,
				Safe, Safe, Safe, Safe, Safe, Safe, Safe, Safe, Safe, Safe,
				Safe, Safe, Safe, Safe, Safe, Safe, Safe, Safe, Safe, Safe,
				Safe,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := genCards(tt.args.count)

			if !cmp.Equal(got, tt.want, cmpopts.SortSlices(
				func(a, b Card) bool {
					return a < b
				},
			)) {
				t.Errorf("genCards() = %v", got)
			}
		})
	}
}

func TestGame_drawCard(t *testing.T) {
	type args struct {
		idx    int
		target int
	}
	tests := []struct {
		g       *Game
		name    string
		args    args
		want    Card
		wantErr bool
	}{
		{
			name: "draw a card",
			g: &Game{
				cards: []Card{
					Bomb, Defuse, Defuse, Defuse, Defuse,
					Safe, Safe, Safe, Safe, Safe,
					Safe, Safe, Safe, Safe, Safe,
					Safe, Safe, Safe, Safe, Safe,
				},
			},
			args: args{
				idx:    2,
				target: 1,
			},
			want: Safe,
		},
		{
			name: "draw a defuse",
			g: &Game{
				cards: []Card{
					Bomb, Defuse, Defuse, Defuse, Defuse,
					Safe, Safe, Safe, Safe, Safe,
					Safe, Safe, Safe, Safe, Safe,
					Safe, Safe, Safe, Safe, Safe,
				},
			},
			args: args{
				idx:    0,
				target: 1,
			},
			want: Defuse,
		},
		{
			name: "out of bound",
			g: &Game{
				cards: []Card{},
			},
			args: args{
				idx:    2,
				target: 1,
			},
			want:    Unkown,
			wantErr: true,
		},
		{
			name: "last round, card from another player",
			g: &Game{
				round: 4,
				cards: []Card{
					Bomb, Defuse, Safe, Safe,
				},
			},
			args: args{
				idx:    0,
				target: 3, // here max should be 0 since everyone has one card
			},
			want:    Unkown,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.drawCard(tt.args.idx, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Game.drawCard() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Game.drawCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGame_nextRoundCards(t *testing.T) {
	type fields struct {
		cards []Card
		ids   []int64
	}
	tests := []struct {
		name      string
		fields    fields
		wantCards []Card
	}{
		{
			name: "round 0 to 1, 20 to 16 cards",
			fields: fields{
				ids: testPlayerIDs,
				cards: []Card{
					Bomb, Drawn, Defuse, Defuse, Defuse,
					Safe, Drawn, Safe, Safe, Safe,
					Safe, Drawn, Safe, Safe, Safe,
					Safe, Drawn, Safe, Safe, Safe,
				},
			},
			wantCards: []Card{
				Bomb, Defuse, Defuse, Defuse,
				Safe, Safe, Safe, Safe,
				Safe, Safe, Safe, Safe,
				Safe, Safe, Safe, Safe,
			},
		},
		{
			name: "round 1 to 2, 16 to 12 cards",
			fields: fields{
				ids: testPlayerIDs,
				cards: []Card{
					Bomb, Drawn, Defuse, Defuse,
					Safe, Drawn, Safe, Safe,
					Safe, Drawn, Safe, Safe,
					Safe, Drawn, Safe, Safe,
				},
			},
			wantCards: []Card{
				Bomb, Defuse, Defuse,
				Safe, Safe, Safe,
				Safe, Safe, Safe,
				Safe, Safe, Safe,
			},
		},
		{
			name: "round 2 to 3, 12 to 8 cards",
			fields: fields{
				ids: testPlayerIDs,
				cards: []Card{
					Bomb, Drawn, Defuse,
					Safe, Drawn, Safe,
					Safe, Drawn, Safe,
					Safe, Drawn, Safe,
				},
			},
			wantCards: []Card{
				Bomb, Defuse,
				Safe, Safe,
				Safe, Safe,
				Safe, Safe,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				cards: tt.fields.cards,
				ids:   tt.fields.ids,
			}
			g.setNextRoundCards()

			if !cmp.Equal(g.cards, tt.wantCards, cmpopts.SortSlices(
				func(a, b Card) bool {
					return a < b
				},
			)) {
				t.Errorf("genCards() = %v", g.cards)
			}
		})
	}
}

func TestGame_getPlayerCards(t *testing.T) {
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
		idx int
	}
	tests := []struct {
		name   string
		want   []Card
		fields fields
		args   args
	}{
		{
			name: "third round, 4 players",
			fields: fields{
				round: 3,
				ids:   testPlayerIDs,
				cards: []Card{
					Bomb, Safe,
					Safe, Safe,
					Safe, Safe,
					Safe, Defuse,
				},
			},
			args: args{
				idx: 3,
			},
			want: []Card{
				Safe, Defuse,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Game{
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

			got := g.getPlayerCards(tt.args.idx)

			if !cmp.Equal(got, tt.want, cmpopts.SortSlices(
				func(a, b Card) bool {
					return a < b
				},
			)) {
				t.Errorf("genCards() = %v", g.cards)
			}
		})
	}
}
