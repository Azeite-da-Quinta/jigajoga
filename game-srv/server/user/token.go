package user

type Identifier string

type Token interface {
	ID() Identifier
	RoomID() Identifier
	Name() string
}
