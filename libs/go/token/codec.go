// Package token wraps JWT code and decode
package token

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// AccessExpiration expected time to connect to svc/game-srv
	AccessExpiration = 10 * time.Minute
	// RefreshExpiration expected duration of a party
	RefreshExpiration = 24 * time.Hour
)

// DefaultSecret in string base64 format.
// THIS MUST BE OVERRIDE. It's weak on purpose so we don't forget
const DefaultSecret = "QWxoZWlyYXM="

var (
	// ErrBadClaims if the cast failed
	ErrBadClaims = errors.New("failed to get claims")
	// ErrMissingContent missing content in token
	ErrMissingContent = errors.New("missing content in token")
)

// Base64ToKey converts a string base64 key into bytes
func Base64ToKey(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// Claims custom claims type for JWT Token.
// Note: by removing the JSON tags, it "hides" the
// envelope type on the final token.
type Claims struct {
	Envelope
	jwt.RegisteredClaims
}

const issuer = "jigajoga-butler"

var audience = []string{"jigajoga-client"}

// Claims creates JWT claims to prepare a token. Use expiration constants
// from this pkg
func (e Envelope) Claims(now time.Time, expiration time.Duration) Claims {
	return Claims{
		Envelope: e,
		// TODO check what fields to use
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: issuer,
			// Subject:   "somebody",
			Audience:  audience,
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			// NotBefore: jwt.NewNumericDate(now), // We don't need this
			IssuedAt: jwt.NewNumericDate(now),
			// ID:        "1",
		},
	}
}

// Codec both encodes and decodes JWT claims
type Codec struct {
	Key []byte
}

// Encode the claims to a signed JWT string
func (c Codec) Encode(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(c.Key)
}

// Decode the JWT string and returns its inner Envelope
func (c Codec) Decode(s string) (Envelope, error) {
	token, err := jwt.ParseWithClaims(
		s, &Claims{},
		func(token *jwt.Token) (any, error) {
			return c.Key, nil
		},
		jwt.WithAudience(audience[0]),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		// TODO: check other options
	)
	if err != nil {
		return Envelope{}, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return Envelope{}, ErrBadClaims
	}

	if claims.Access == nil &&
		claims.Refresh == nil {
		return Envelope{}, ErrMissingContent
	}

	return claims.Envelope, nil
}
