package bomb

import "fmt"

// ErrCode type
type ErrCode uint8

// enum kinds of error codes
const (
	Nothing ErrCode = iota
	MinPlayers
	PartyNotReady
	PartyAlreadyStarted
	TooManyPlayers
	PlayerNotFound
	TargetNotFound
	NotYourTurn
	CardNotFound
	CannotDrawSelf // cannot draw from your own hand
	// TODO ^
)

// game errors
var (
	ErrMinPlayers = Error{
		reason:  "not enough players to start the game",
		ErrCode: MinPlayers,
	}
	ErrPartyNotReady = Error{
		reason:  "some party players are not ready",
		ErrCode: PartyNotReady,
	}
	ErrPartyAlreadyStarted = Error{
		reason:  "party already started",
		ErrCode: PartyAlreadyStarted,
	}
	ErrTooManyPlayers = Error{
		reason:  "party is already full",
		ErrCode: TooManyPlayers,
	}
	ErrPlayerNotFound = Error{
		reason:  "player not found",
		ErrCode: PlayerNotFound,
	}
	ErrTargetNotFound = Error{
		reason:  "player target not found",
		ErrCode: TargetNotFound,
	}
	ErrNotYourTurn = Error{
		reason:  "not your turn",
		ErrCode: NotYourTurn,
	}
	ErrCardNotFound = Error{
		reason:  "card not found",
		ErrCode: CardNotFound,
	}
)

// Error type implements error and adds a code information
// to be used by a client/frontend. Useful for i18n
type Error struct {
	reason string
	ErrCode
}

func (e Error) Error() string {
	return fmt.Sprintf("code=%d reason=%s ", e.ErrCode, e.reason)
}
