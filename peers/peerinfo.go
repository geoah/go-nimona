package peers

import (
	"nimona.io/go/encoding"
)

// PeerInfo holds the information exchange needs to connect to a remote peer
type PeerInfo struct {
	Addresses []string         `json:"addresses"`
	RawObject *encoding.Object `json:"@"`
}

func (pi *PeerInfo) Thumbprint() string {
	// TODO(geoah) should this return the authority or the subject's id?
	return pi.RawObject.HashBase58()
}

// Address of the peer
func (pi *PeerInfo) Address() string {
	return "peer:" + pi.RawObject.SignerKey().HashBase58()
}
