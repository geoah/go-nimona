package net

import (
	"nimona.io/go/crypto"
)

type HandshakeAck struct {
	Nonce     string            `json:"nonce"`
	Signature *crypto.Signature `json:"@sig"`
}
