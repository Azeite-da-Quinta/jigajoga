package bomb

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

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
