package token

import "time"

// Envelope contains either access or refresh fields
type Envelope struct {
	*Access  `json:"access,omitempty"`
	*Refresh `json:"refresh,omitempty"`
}

func (e Envelope) expiration() time.Duration {
	if e.Access != nil {
		return AccessExpiration
	} else if e.Refresh != nil {
		return RefreshExpiration
	}

	return 0
}
