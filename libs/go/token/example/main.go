// example usage of token
package main

import (
	"fmt"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/libs/go/token"
)

// go run libs/go/token/example/main.go
func main() {
	b, err := token.Base64ToKey(token.DefaultSecret)
	if err != nil {
		panic(err)
	}

	codec := token.Codec{
		Key: b,
	}

	e := token.Envelope{
		Access: &token.Access{
			RoomIDField: "1234",
			IDField:     "1",
			NameField:   "bob",
		},
	}

	s, err := codec.Encode(e.Claims(time.Now()))
	if err != nil {
		panic(err)
	}

	fmt.Println("access token:")
	fmt.Println(s)

	e = token.Envelope{
		Refresh: &token.Refresh{
			PlayerIDField: "001",
			Rotation:      36,
		},
	}

	s, err = codec.Encode(e.Claims(time.Now()))
	if err != nil {
		panic(err)
	}

	fmt.Println("refresh token:")
	fmt.Println(s)
}
