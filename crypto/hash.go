package crypto

import (
	"golang.org/x/crypto/sha3"
	"nimona.io/go/base58"
	"nimona.io/go/encoding"
)

//proteus:generate
type Sha3 struct {
	Type string   `json:"@ctx"`
	Hash [32]byte `json:"hash"`
}

func (h *Sha3) Base58() string {
	b, err := encoding.Marshal(h)
	if err != nil {
		panic(err)
	}

	return base58.Encode(b)
}

func NewSha3(b []byte) *Sha3 {
	return &Sha3{
		Type: "sha3.256",
		Hash: sha3.Sum256(b),
	}
}
