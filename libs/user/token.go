// Package user for infos about active users on the server
package user

import (
	"strconv"

	"github.com/Azeite-da-Quinta/jigajoga/libs/token"
)

// Identifier type for the app
type Identifier int64

// Token represents the identity of a User
type Token interface {
	ID() Identifier
	RoomID() Identifier
	Name() string
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

// ToToken converts
func (d Data) ToToken() (t token.Data) {
	const base = 10

	return token.Data{
		IDField:     strconv.FormatInt(int64(d.id), base),
		RoomIDField: strconv.FormatInt(int64(d.room), base),
		NameField:   d.name,
	}
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
