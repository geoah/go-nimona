package net

import (
	"nimona.io/go/crypto"
)

// ForwardRequest is the payload for proxied blocks
type ForwardRequest struct {
	Recipient string            `json:"recipient"` // address
	FwBlock   interface{}       `json:"fwBlock"`
	Signature *crypto.Signature `json:"@signature"`
}
