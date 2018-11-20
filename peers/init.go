package peers

import (
	"nimona.io/go/encoding"
)

func init() {
	encoding.Register("/peer.info", &PeerInfo{})
}
