package primitives

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	ucodec "github.com/ugorji/go/codec"

	"github.com/pkg/errors"
	"nimona.io/go/base58"
)

// Supported values for KeyType
const (
	EC             = "EC"  // Elliptic Curve
	InvalidKeyType = ""    // Invalid KeyType
	OctetSeq       = "oct" // Octet sequence (used to represent symmetric keys)
	RSA            = "RSA" // RSA
)

const (
	// P256 curve
	P256 string = "P-256"
	// P384 curve
	P384 string = "P-384"
	// P521 curve
	P521 string = "P-521"
)

// Key defines the minimal interface for each of the
// key types.
type Key struct {
	Algorithm              string `json:"alg,omitempty" mapstructure:"alg,omitempty"`
	KeyID                  string `json:"kid,omitempty" mapstructure:"kid,omitempty"`
	KeyType                string `json:"kty,omitempty" mapstructure:"kty,omitempty"`
	KeyUsage               string `json:"use,omitempty" mapstructure:"use,omitempty"`
	KeyOps                 string `json:"key_ops,omitempty" mapstructure:"key_ops,omitempty"`
	X509CertChain          string `json:"x5c,omitempty" mapstructure:"x5c,omitempty"`
	X509CertThumbprint     string `json:"x5t,omitempty" mapstructure:"x5t,omitempty"`
	X509CertThumbprintS256 string `json:"x5tS256,omitempty" mapstructure:"x5tS256,omitempty"`
	X509URL                string `json:"x5u,omitempty" mapstructure:"x5u,omitempty"`
	Curve                  string `json:"crv,omitempty" mapstructure:"crv,omitempty"`
	X                      []byte `json:"x,omitempty" mapstructure:"x,omitempty"`
	Y                      []byte `json:"y,omitempty" mapstructure:"y,omitempty"`
	D                      []byte `json:"d,omitempty" mapstructure:"d,omitempty"`
	key                    interface{}
}

func (b *Key) Block() *Block {
	s := structs.New(b)
	s.TagName = "mapstructure"
	return &Block{
		Type:    "nimona.io/key",
		Payload: s.Map(),
	}
}

func (b *Key) FromBlock(block *Block) {
	mapstructure.Decode(block.Payload, b)
}

// CodecDecodeSelf helper for cbor unmarshaling
func (k *Key) CodecDecodeSelf(dec *ucodec.Decoder) {
	b := &Block{}
	dec.MustDecode(b)
	k.FromBlock(b)
}

// CodecEncodeSelf helper for cbor marshaling
func (k *Key) CodecEncodeSelf(enc *ucodec.Encoder) {
	b := k.Block()
	enc.MustEncode(b)
}

func (k *Key) Thumbprint() string {
	b, _ := Marshal(k.Block())
	return base58.Encode(b)
}

// GetPublicKey returns the public key
func (k *Key) GetPublicKey() *Key {
	if len(k.D) == 0 {
		return k
	}

	pk := k.Materialize().(*ecdsa.PrivateKey).Public().(*ecdsa.PublicKey)
	bpk, err := NewKey(pk)
	if err != nil {
		panic(err)
	}

	return bpk
}

func (k *Key) Materialize() interface{} {
	// TODO cache on k.key
	var curve elliptic.Curve
	switch k.Curve {
	case P256:
		curve = elliptic.P256()
	case P384:
		curve = elliptic.P384()
	case P521:
		curve = elliptic.P521()
	default:
		panic("invalid curve name " + k.Curve)
		// return nil, errors.Errorf(`invalid curve name %s`, h.Curve)
	}

	var key interface{}
	switch k.KeyType {
	case EC:
		if len(k.D) > 0 {
			key = &ecdsa.PrivateKey{
				PublicKey: ecdsa.PublicKey{
					Curve: curve,
					X:     bigIntFromBytes(k.X),
					Y:     bigIntFromBytes(k.Y),
				},
				D: bigIntFromBytes(k.D),
			}
		} else {
			key = &ecdsa.PublicKey{
				Curve: curve,
				X:     bigIntFromBytes(k.X),
				Y:     bigIntFromBytes(k.Y),
			}
		}
	default:
		panic("invalid kty")
		// return nil, errors.Errorf(`invalid kty %s`, h.KeyType)
	}

	return key
}

// NewKey creates a Key from the given key.
func NewKey(k interface{}) (*Key, error) {
	if k == nil {
		return nil, errors.New("missing key")
	}

	key := &Key{
		key: k,
	}

	switch v := k.(type) {
	// case *rsa.PrivateKey:
	// 	return newRSAPrivateKey(v)
	// case *rsa.PublicKey:
	// 	return newRSAPublicKey(v)
	case *ecdsa.PrivateKey:
		key.KeyType = EC
		key.Curve = v.Curve.Params().Name
		key.X = v.X.Bytes()
		key.Y = v.Y.Bytes()
		key.D = v.D.Bytes()
	case *ecdsa.PublicKey:
		key.KeyType = EC
		key.Curve = v.Curve.Params().Name
		key.X = v.X.Bytes()
		key.Y = v.Y.Bytes()
	// case []byte:
	// 	return newSymmetricKey(v)
	default:
		return nil, errors.Errorf(`invalid key type %T`, key)
	}

	return key, nil
}

func bigIntFromBytes(b []byte) *big.Int {
	i := &big.Int{}
	return i.SetBytes(b)
}
