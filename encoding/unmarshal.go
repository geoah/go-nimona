package encoding

import (
	"github.com/ugorji/go/codec"
)

func Unmarshal(b []byte) (interface{}, error) {
	m := map[string]interface{}{}
	dec := codec.NewDecoderBytes(b, CborHandler())
	if err := dec.Decode(m); err != nil {
		return nil, err
	}

	v := GetInstance(GetType(m))
	if v == nil {
		v = &map[string]interface{}{}
	}
	if err := Decode(m, v, true); err != nil {
		return nil, err
	}

	return v, nil
}

func UnmarshalInto(b []byte, v interface{}) error {
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
