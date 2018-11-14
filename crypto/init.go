package crypto

import (
	"nimona.io/go/encoding"
)

func init() {
	encoding.Register("/key", &Key{})
	encoding.Register("/sig", &Signature{})

	encoding.Register("/policy", &Policy{})
}
