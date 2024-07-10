// example usage of token
package main

import (
	"fmt"
	"time"

	"github.com/Azeite-da-Quinta/jigajoga/libs/go/token"
)

func main() {
	b, err := token.Base64ToKey(token.DefaultSecret)
	if err != nil {
		panic(err)
	}

	codec := token.Codec{
		Key: b,
	}

	d := token.Data{
		RoomIDField: "1234",
		IDField:     "1",
		NameField:   "bob",
	}

	s, err := codec.Encode(d.Claims(time.Now()))
	if err != nil {
		panic(err)
	}

	fmt.Println(s)
}
