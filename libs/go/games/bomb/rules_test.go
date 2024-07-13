package bomb

import (
	"testing"
)

func Test_isRoundEnd(t *testing.T) {
	type args struct {
		playersCount int
		cardsDealt   int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "revealed one card",
			args: args{
				playersCount: 4,
				cardsDealt:   1,
			},
			want: false,
		},
		{
			name: "revealed last card of the round",
			args: args{
				playersCount: 4,
				cardsDealt:   4,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRoundEnd(tt.args.playersCount, tt.args.cardsDealt); got != tt.want {
				t.Errorf("isRoundEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_areRoundsOver(t *testing.T) {
	type args struct {
		round uint8
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "first round",
			args: args{},
			want: false,
		},
		{
			name: "last round",
			args: args{
				round: 3,
			},
			want: false,
		},
		{
			name: "passed last round",
			args: args{
				round: 4,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := areRoundsOver(tt.args.round); got != tt.want {
				t.Errorf("areRoundsOver() = %v, want %v", got, tt.want)
			}
		})
	}
}
