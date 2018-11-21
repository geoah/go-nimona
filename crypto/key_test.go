package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"nimona.io/go/base58"
	"nimona.io/go/encoding"
)

func TestKeyEncoding(t *testing.T) {
	ek := &Key{
		Algorithm: "key-alg",
	}

	em := map[string]interface{}{
		"@ctx": "/key",
		"alg":  "key-alg",
	}

	bs, err := encoding.Marshal(ek)
	assert.NoError(t, err)

	assert.Equal(t, "Nx3cnuT6J8XPNCBmncEt5BfwfKtYtf6h5VzoQ", base58.Encode(bs))

	m := map[string]interface{}{}
	err = encoding.UnmarshalInto(bs, &m)

	assert.Equal(t, em, m)

	k := &Key{}
	err = encoding.UnmarshalInto(bs, k)

	assert.Equal(t, ek, k)
}
