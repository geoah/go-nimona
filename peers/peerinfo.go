package peers

import (
	"nimona.io/go/crypto"
	"nimona.io/go/encoding"
)

// PeerInfo holds the information exchange needs to connect to a remote peer
type PeerInfo struct {
	Addresses []string            `json:"addresses"`
	Signature *crypto.Signature   `json:"@sig"`
	Raw       *encoding.Container `json:"@"`
}

func (pi *PeerInfo) Thumbprint() string {
	return pi.Signature.Key.Thumbprint()
}

func (pi *PeerInfo) Address() string {
	return "peer:" + pi.Signature.Key.Thumbprint()
}
