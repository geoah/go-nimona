package net

import (
	"nimona.io/go/crypto"
	"nimona.io/go/peers"
)

type HandshakeSynAck struct {
	Nonce     string            `json:"nonce"`
	PeerInfo  *peers.PeerInfo   `json:"peerInfo,omitempty"`
	Signature *crypto.Signature `json:"@sig"`
}
