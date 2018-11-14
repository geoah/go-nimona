package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"nimona.io/go/base58"
	"nimona.io/go/encoding"
)

// func TestSignatureVerification(t *testing.T) {
// 	sk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, sk)

// 	k, err := NewKey(sk)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, sk)

// 	p := &Block{
// 		Type: "test.nimona.io/dummy",
// 		Payload: map[string]interface{}{
// 			"foo": "bar2",
// 		},
// 	}

// 	err = Sign(p, k)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, p.Signature)

// 	// test verification
// 	digest, err := getDigest(p)
// 	assert.NoError(t, err)
// 	err = Verify(p.Signature, digest)
// 	assert.NoError(t, err)
// }

func TestSignatureEncoding(t *testing.T) {
	es := &Signature{
		Key: &Key{
			Algorithm: "key-alg",
		},
		Alg: "sig-alg",
	}

	em := map[string]interface{}{
		"@ctx": "/sig",
		"alg":  "sig-alg",
		"key":  es.Key,
	}

	bs, err := encoding.Marshal(es)
	assert.NoError(t, err)

	assert.Equal(t, "BiB8mh2HL54SZ82m9tEsbMRGPC7QML6jkCc2G2YTfvkW2kBaC3D9XZNdw"+
		"aqYPvTMrkGCxQGNrqJQytWng83bsDYLeb9xerPVYhdP", base58.Encode(bs))

	m := map[string]interface{}{}
	err = encoding.UnmarshalInto(bs, &m)

	assert.Equal(t, em, m)

	s := &Signature{}
	err = encoding.UnmarshalInto(bs, s)
	assert.Equal(t, es, s)
}
