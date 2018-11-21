package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
)

var (
	// ErrInvalidBlockType is returned when the signature being verified
	// is not an encoded block of type "signature".
	ErrInvalidBlockType = errors.New("invalid block type")
	// ErrAlgorithNotImplemented is returned when the algorithm specified
	// has not been implemented
	ErrAlgorithNotImplemented = errors.New("algorithm not implemented")
)

// Signature block (container), currently supports only ES256
type Signature struct {
	Alg string `json:"alg"`
	R   []byte `json:"r"`
	S   []byte `json:"s"`
}

// NewSignature returns a signature given some bytes and a private key
func NewSignature(key *Key, alg string, digest []byte) (*Signature, error) {
	if key == nil {
		return nil, errors.New("missing key")
	}

	mKey := key.Materialize()
	if mKey == nil {
		return nil, errors.New("could not materialize")
	}

	pKey, ok := mKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("only ecdsa private keys are currently supported")
	}

	if alg != "ES256" {
		return nil, ErrAlgorithNotImplemented
	}

	// TODO implement more algorithms
	hash := sha256.Sum256(digest)
	r, s, err := ecdsa.Sign(rand.Reader, pKey, hash[:])
	if err != nil {
		return nil, err
	}

	return &Signature{
		Alg: alg,
		R:   r.Bytes(),
		S:   s.Bytes(),
	}, nil
}
