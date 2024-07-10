// Package slogt contains common slog types
package slogt

import "log/slog"

// Error common error attribute
func Error(err error) slog.Attr {
	return slog.Attr{Key: "err", Value: slog.StringValue(err.Error())}
}

// PlayerID common ID attribute
func PlayerID(id int64) slog.Attr {
	return slog.Attr{Key: "player_id", Value: slog.Int64Value(id)}
}

// PlayerName common name attribute
func PlayerName(name string) slog.Attr {
	return slog.Attr{Key: "player_name", Value: slog.StringValue(name)}
}

// Room common room's ID attribute
func Room(id int64) slog.Attr {
	return slog.Attr{Key: "room_id", Value: slog.Int64Value(id)}
}

// Num of a worker
func Num(num int) slog.Attr {
	return slog.Attr{Key: "num", Value: slog.IntValue(num)}
}
