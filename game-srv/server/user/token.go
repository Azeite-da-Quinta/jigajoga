// Package user for infos about active users on the server
package user

// Identifier for IDs. TODO: replace by UUID maybe?
type Identifier string

// Token represents the identity of a User
type Token interface {
	ID() Identifier
	RoomID() Identifier
	Name() string
}
