package token

import (
	"errors"
	"testing"
	"time"
)

func TestCodec_Encode_Decode(t *testing.T) {
	val := "test value"

	c := Codec{
		Key: []byte("test_key"),
	}

	t.Run("basic case", func(t *testing.T) {
		claims := Envelope{Access: &Access{IDField: val}}.Claims(time.Now(), AccessExpiration)

		s, err := c.Encode(claims)
		if err != nil {
			t.Error("failed to encode", err)
		}

		u, err := c.Decode(s)
		if err != nil {
			t.Error("failed to decode", err)
		}

		if val != u.IDField {
			t.Error("value does not match")
		}
	})

	t.Run("missing content", func(t *testing.T) {
		claims := Envelope{}.Claims(time.Now(), AccessExpiration)

		s, err := c.Encode(claims)
		if err != nil {
			t.Error("failed to encode", err)
		}

		_, err = c.Decode(s)
		if !errors.Is(err, ErrMissingContent) {
			t.Error("decode unexpected error", err)
		}
	})
}
