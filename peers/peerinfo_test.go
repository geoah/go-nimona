package peers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"nimona.io/go/base58"
	"nimona.io/go/crypto"
	"nimona.io/go/encoding"
)

func TestPeerInfoBlock(t *testing.T) {
	ep := &PeerInfo{
		Addresses: []string{
			"p1-addr1",
			"p1-addr2",
		},
		Signature: &crypto.Signature{
			Key: &crypto.Key{
				Algorithm: "key-alg",
			},
			Alg: "sig-alg",
		},
	}

	b := ep
	bs, _ := encoding.Marshal(b)

	b2, err := encoding.Unmarshal(bs)
	assert.NoError(t, err)

	p := &PeerInfo{}
	p.FromBlock(b2)

	assert.Equal(t, ep, p)
}

func TestPeerInfoSelfEncode(t *testing.T) {
	eb := &PeerInfo{
		Addresses: []string{
			"p1-addr1",
			"p1-addr2",
		},
		Signature: &crypto.Signature{
			Key: &crypto.Key{
				Algorithm: "key-alg",
			},
			Alg: "sig-alg",
		},
	}

	bs, err := encoding.Marshal(eb)
	assert.NoError(t, err)

	assert.Equal(t, base58.Encode(bs), "BvE6Qe57DKXhLXzNVg4HeDf6Gv3jFAmZzixdtB"+
		"jLmkQQP9tBLgsRtCPRDF5gUnt4FXZuxMNNbJwScDHwgnRr1SZyK9fNv7zUpV2LyQPCGXk"+
		"wDQ5rGruw8bTfjvfyg9gQiPTWEH5JtCNocVJiEAqj9qrFYtmrKsVDibAL5EJ53dxZrb5M"+
		"UPArU2ze2Yy7jhpib1YxGNZv89WAACh9E4fRRbDQmWaoMf2BLiS")

	b := &PeerInfo{}
	err = encoding.Unmarshal(bs, b)
	assert.NoError(t, err)

	assert.Equal(t, eb, b)
}
