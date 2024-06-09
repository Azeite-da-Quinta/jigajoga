package user

import "github.com/Azeite-da-Quinta/jigajoga/libs/flakes"

// Factory generates users with adequate Identifiers
type Factory struct {
	gen flakes.Generator
}

// NewFactory requires the node's number
func NewFactory(node int64) (Factory, error) {
	gen, err := flakes.New(node)
	if err != nil {
		return Factory{}, err
	}

	return Factory{
		gen: gen,
	}, nil
}

// NewUser generates a User in that room
func (f Factory) NewUser(name string, room Identifier) Data {
	return Data{
		name: name,
		id:   Identifier(f.gen.ID()),
		room: room,
	}
}

// NewRoom generates an ID
func (f Factory) NewRoom() Identifier {
	return Identifier(f.gen.ID())
}

// NewFirstUser TODO maybe not needed
func (f Factory) NewFirstUser(name string) Data {
	return Data{
		name: name,
		id:   Identifier(f.gen.ID()),
		room: Identifier(f.gen.ID()),
	}
}
