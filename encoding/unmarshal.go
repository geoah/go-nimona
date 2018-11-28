package encoding

import (
	"github.com/ugorji/go/codec"
)

// Unmarshal a cbor encoded block (container) into a registered type, or map
func Unmarshal(b []byte) (*Object, error) {
	m := map[string]interface{}{}
	dec := codec.NewDecoderBytes(b, CborHandler())
	if err := dec.Decode(m); err != nil {
		return nil, err
	}

	return NewObject(m), nil
}
