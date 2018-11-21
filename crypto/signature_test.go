package crypto

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"nimona.io/go/base58"
	"nimona.io/go/encoding"
)

func TestSignatureEncoding(t *testing.T) {
	es := &Signature{
		Alg: "sig-alg",
		R:   big.NewInt(math.MaxInt64).Bytes(),
		S:   big.NewInt(math.MinInt64).Bytes(),
	}

	em := map[string]interface{}{
		"@ctx": "/sig",
		"alg":  "sig-alg",
		"r":    es.R,
		"s":    es.S,
	}

	bs, err := encoding.Marshal(es)
	assert.NoError(t, err)

	assert.Equal(t, "41BGbraog8gf47JJunL1pYsE8eeE11v6uirNPoFNuKbJ7tUATDbuXpM2G"+
		"ZqbbGfojVkcPzYfZ", base58.Encode(bs))

	m := map[string]interface{}{}
	err = encoding.UnmarshalInto(bs, &m)

	assert.Equal(t, em, m)

	s := &Signature{}
	err = encoding.UnmarshalInto(bs, s)
	assert.Equal(t, es, s)
}
