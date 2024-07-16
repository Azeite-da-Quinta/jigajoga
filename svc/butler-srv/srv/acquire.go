package srv

import (
	"encoding/json"
	"github.com/Azeite-da-Quinta/jigajoga/libs/go/token"
	"net/http"
	"time"
)

func AcquireHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var user token.Data
		defer r.Body.Close()
		err := decoder.Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			//w.Write([]byte(err.Error()))
		} else {
			b, err := token.Base64ToKey(token.DefaultSecret)
			if err != nil {
				panic(err)
			}

			codec := token.Codec{
				Key: b,
			}

			s, err := codec.Encode(user.Claims(time.Now()))
			if err != nil {
				panic(err)
			}

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(s))
		}

	}
}
