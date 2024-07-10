// Package user for infos about active users on the server
package user

import (
	"encoding/json"
	"strconv"

	"github.com/Azeite-da-Quinta/jigajoga/libs/go/token"
)

// Identifier type for the app
type Identifier int64

// Token represents the identity of a User
type Token interface {
	ID() Identifier
	RoomID() Identifier
	Name() string
	Serialize() ([]byte, error)
}

// proof of implementation
var _ Token = Data{}

// Data of a user
type Data struct {
	name string
	id   Identifier
	room Identifier
}

// FromToken creates a User from the token's data
func FromToken(t token.Data) (Data, error) {
	var u Data

	id, err := strconv.Atoi(t.IDField)
	if err != nil {
		return u, err
	}

	room, err := strconv.Atoi(t.RoomIDField)
	if err != nil {
		return u, err
	}

	u.id = Identifier(id)
	u.room = Identifier(room)
	u.name = t.NameField

	return u, nil
}

// Token converts
func (d Data) Token() (t token.Data) {
	const base = 10

	return token.Data{
		IDField:     strconv.FormatInt(int64(d.id), base),
		RoomIDField: strconv.FormatInt(int64(d.room), base),
		NameField:   d.name,
	}
}

// Deserialize returns
func Deserialize(data []byte) (Token, error) {
	var t token.Data
	err := json.Unmarshal(data, &t)
	if err != nil {
		return Data{}, err
	}

	return FromToken(t)
}

// Serialize returns json bytes
func (d Data) Serialize() ([]byte, error) {
	return json.Marshal(d.Token())
}

// ID implements Token
func (d Data) ID() Identifier {
	return d.id
}

// RoomID implements Token
func (d Data) RoomID() Identifier {
	return d.room
}

// Name implements Token
func (d Data) Name() string {
	return d.name
}
