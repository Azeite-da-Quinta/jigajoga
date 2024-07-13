package bomb

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

//revive:disable:add-constant
func Test_genRoles(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name string
		want []Role
		args args
	}{
		{
			name: "empty",
			args: args{
				count: 0,
			},
			want: []Role{},
		},
		{
			name: "small group-",
			args: args{
				count: 4,
			},
			want: []Role{Hero, Hero, Vilain, Vilain},
		},
		{
			name: "small group",
			args: args{
				count: 5,
			},
			want: []Role{Hero, Hero, Hero, Vilain, Vilain},
		},
		{
			name: "mid group",
			args: args{
				count: 6,
			},
			want: []Role{Hero, Hero, Hero, Hero, Vilain, Vilain},
		},
		{
			name: "big group",
			args: args{
				count: 8,
			},
			want: []Role{Hero, Hero, Hero, Hero, Hero, Vilain, Vilain, Vilain},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := genRoles(tt.args.count)

			if len(got) != len(tt.want) {
				t.Errorf("genRoles() size = %v, want %v", len(got), len(tt.want))
			}

			if len(got) == 4 {
				return
			}

			if !cmp.Equal(got, tt.want, cmpopts.SortSlices(
				func(a, b Role) bool {
					return a < b
				},
			)) {
				t.Errorf("genRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}
