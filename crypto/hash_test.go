package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"nimona.io/go/base58"
)

func TestHashSha3(t *testing.T) {
	p := []byte{0x7c, 0xc3, 0x6f, 0x04, 0x0a, 0x82, 0x47, 0xbb}
	h := NewSha3(p)
	assert.NotNil(t, h)
	assert.Equal(t, "SHA3", h.Alg)
	assert.NotEmpty(t, h.Digest)

	b, err := base58.Decode(h.Base58())
	assert.NoError(t, err)
	assert.Equal(t, h.Bytes(), b)
}
