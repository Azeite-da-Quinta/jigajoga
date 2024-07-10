// Package kvstore around Redis
package kvstore

import "testing"

func Test_getEvent(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "passing case",
			args: args{
				s: "rooms.1251986239494553600.join",
			},
			want: "join",
		},
		{
			name: "wrong size",
			args: args{
				s: "rooms.join",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEvent(tt.args.s); got != tt.want {
				t.Errorf("getEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
