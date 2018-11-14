package dht

import (
	"nimona.io/go/encoding"
)

func init() {
	encoding.Register("/dht.peer.request", &PeerInfoRequest{})
	encoding.Register("/dht.peer.response", &PeerInfoResponse{})

	encoding.Register("/dht.provider.response", &ProviderRequest{})
	encoding.Register("/dht.provider.response", &ProviderResponse{})
}
