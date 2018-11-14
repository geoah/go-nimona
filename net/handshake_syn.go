package net

import (
	"nimona.io/go/crypto"
	"nimona.io/go/peers"
)

//proteus:generate
type HandshakeSyn struct {
	Nonce     string            `json:"nonce"`
	PeerInfo  *peers.PeerInfo   `json:"peerInfo,omitempty"`
	Signature *crypto.Signature `json:"@sig"`
}
