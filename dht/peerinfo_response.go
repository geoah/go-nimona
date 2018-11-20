package dht

import (
	"nimona.io/go/crypto"
	"nimona.io/go/peers"
)

type PeerInfoResponse struct {
	RequestID    string            `json:"requestID,omitempty"`
	PeerInfo     *peers.PeerInfo   `json:"peerInfo,omitempty"`
	ClosestPeers []*peers.PeerInfo `json:"closestPeers,omitempty"`
	Signature    *crypto.Signature `json:"@sig"`
}
