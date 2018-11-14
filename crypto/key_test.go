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

	assert.Equal(t, "BvuiBQ7bekExgi62fxZ2ARcqoCPYY3P5xQjc2hwiFaSS", base58.Encode(bs))

	m := map[string]interface{}{}
	err = encoding.UnmarshalInto(bs, &m)

	assert.Equal(t, em, m)

	k := &Key{}
	err = encoding.UnmarshalInto(bs, k)

	assert.Equal(t, ek, k)
}
