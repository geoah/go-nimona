package dht

import (
	"nimona.io/go/crypto"
)

// PeerInfoRequest payload
//proteus:generate
type PeerInfoRequest struct {
	RequestID string            `json:"requestID,omitempty"`
	PeerID    string            `json:"peerID"`
	Signature *crypto.Signature `json:"@sig"`
}
