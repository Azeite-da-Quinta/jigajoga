package token

// Envelope contains either access or refresh fields
type Envelope struct {
	*Access  `json:"access,omitempty"`
	*Refresh `json:"refresh,omitempty"`
}
