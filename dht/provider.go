package dht

import (
	"nimona.io/go/crypto"
)

// Provider payload
//proteus:generate
type Provider struct {
	BlockIDs  []string          `json:"blockIDs"`
	Signature *crypto.Signature `json:"@sig"`
}
