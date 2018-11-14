package dht

import (
	"nimona.io/go/crypto"
	"nimona.io/go/peers"
)

//proteus:generate
type ProviderResponse struct {
	RequestID    string            `json:"requestID,omitempty"`
	Providers    []*Provider       `json:"providers,omitempty"`
	ClosestPeers []*peers.PeerInfo `json:"closestPeers,omitempty"`
	Signature    *crypto.Signature `json:"@sig"`
}
