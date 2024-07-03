// Package token wraps JWT code and decode
package token

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//revive:disable:line-length-limit

// DefaultSecret in string base64 format
const DefaultSecret = "QSBhbGhlaXJhIMOpIHVtIGVuY2hpZG8gdMOtcGljbyBkYSBjdWxpbsOhcmlhIHBvcnR1Z3Vlc2EgY3Vqb3MgcHJpbmNpcGFpcyBpbmdyZWRpZW50ZXMgc8OjbyBjYXJuZSBkZSBhdmVzLCBww6NvLCBhemVpdGUsIGJhbmhhLCBhbGhvIGUgY29sb3JhdS4="

//revive:enable:line-length-limit

// ErrBadClaims if the cast failed
var ErrBadClaims = errors.New("failed to get claims")

// Base64ToKey converts a string base64 key into bytes
func Base64ToKey(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

// Claims custom claims type for JWT Token
type Claims struct {
	Data `json:"u"`
	jwt.RegisteredClaims
}

// Claims creates JWT claims to prepare a token
func (u Data) Claims(now time.Time) Claims {
	return Claims{
		Data: u,
		// TODO check what fields to use
		RegisteredClaims: jwt.RegisteredClaims{
			// A usual scenario is to set the expiration
			// time relative to the current time
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			// IssuedAt:  jwt.NewNumericDate(now),
			// NotBefore: jwt.NewNumericDate(now),
			Issuer: "jigajoga",
			// Subject:   "somebody",
			// ID:        "1",
			// Audience:  []string{"somebody_else"},
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

// Decode the JWT string and returns its inner Data
func (c Codec) Decode(s string) (Data, error) {
	token, err := jwt.ParseWithClaims(s, &Claims{},
		func(token *jwt.Token) (any, error) {
			return c.Key, nil
		},
	)

	if err != nil {
		return Data{}, err
	} else if claims, ok := token.Claims.(*Claims); ok {
		return claims.Data, nil
	}

	return Data{}, ErrBadClaims
}
