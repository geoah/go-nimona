package encoding

import (
	"github.com/ugorji/go/codec"
)

func UnmarshalSimple(b []byte, v interface{}) error {
	dec := codec.NewDecoderBytes(b, CborHandler())
	return dec.Decode(v)
}
