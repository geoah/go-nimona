package dht

import (
	"nimona.io/go/crypto"
)

// ProviderRequest payload
//proteus:generate
type ProviderRequest struct {
	RequestID string            `json:"requestID,omitempty"`
	Key       string            `json:"key"`
	Signature *crypto.Signature `json:"@sig"`
}
