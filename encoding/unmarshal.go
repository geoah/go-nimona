package encoding

import (
	"github.com/ugorji/go/codec"
)

// Unmarshal a cbor encoded block (container) into a registered type, or map
func Unmarshal(b []byte) (interface{}, error) {
	m := map[string]interface{}{}
	dec := codec.NewDecoderBytes(b, CborHandler())
	if err := dec.Decode(m); err != nil {
		return nil, err
	}

	um, err := UntypeMap(m)
	if err != nil {
		return nil, err
	}

	v := GetInstance(GetType(um))
	if v == nil {
		v = &map[string]interface{}{}
	}
	if err := Decode(um, v, true); err != nil {
		return nil, err
	}

	return v, nil
}

// UnmarshalInto unmarshals a cbor encoded block (container) into a given type
func UnmarshalInto(b []byte, v interface{}) error {
	m := map[string]interface{}{}
	dec := codec.NewDecoderBytes(b, CborHandler())
	if err := dec.Decode(m); err != nil {
		return err
	}

	um, err := UntypeMap(m)
	if err != nil {
		return err
	}

	if err := Decode(um, v, true); err != nil {
		return err
	}

	return nil
}
