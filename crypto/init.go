package crypto

import (
	"nimona.io/go/encoding"
)

func init() {
	encoding.Register("/key", &Key{})
	encoding.Register("/sig", &Signature{})
	encoding.Register("/mandate", &Mandate{})
	encoding.Register("/mandate.policy", &MandatePolicy{})
	encoding.Register("/policy", &Policy{})
}
