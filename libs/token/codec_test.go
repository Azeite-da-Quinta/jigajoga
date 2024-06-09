package token

import (
	"testing"
	"time"
)

func TestCodec_Encode_Decode(t *testing.T) {
	t.Run("basic case", func(t *testing.T) {
		val := "test value"

		now := time.Now()

		claims := Data{IDField: val}.Claims(now)

		c := Codec{
			Key: []byte("test_key"),
		}

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
}
