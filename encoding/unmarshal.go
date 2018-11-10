package encoding

import (
	"github.com/ugorji/go/codec"
)

func Unmarshal(b []byte, v interface{}) error {
	m := map[string]interface{}{}
	dec := codec.NewDecoderBytes(b, CborHandler())
	if err := dec.Decode(m); err != nil {
		return err
	}

	if err := Decode(m, v, true); err != nil {
		return err
	}

	return nil
}
