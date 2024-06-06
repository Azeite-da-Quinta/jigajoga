// Package slogt contains common slog types
package slogt

import "log/slog"

// Error common error attribute
func Error(err error) slog.Attr {
	return slog.Attr{Key: "err", Value: slog.StringValue(err.Error())}
}

// ID common ID attribute
func ID(id int) slog.Attr {
	return slog.Attr{Key: "id", Value: slog.IntValue(id)}
}
