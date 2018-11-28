package encoding

import (
	"github.com/ugorji/go/codec"
)

// Marshal encodes an object to cbor
func Marshal(o *Object) ([]byte, error) {
	m := o.Map()
	b := []byte{}
	enc := codec.NewEncoderBytes(&b, CborHandler())
	err := enc.Encode(m)
	return b, err
}
