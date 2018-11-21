package encoding

import (
	"errors"
	"reflect"

	"github.com/ugorji/go/codec"
)

var (
	// ErrUnknownType for when we are trying to marshal an unregistered type
	ErrUnknownType = errors.New("unknown type")
)

// Marshal encodes a given struct to a block (container), and encodes it as cbor
func Marshal(v interface{}) ([]byte, error) {
	rt := GetType(reflect.TypeOf(v))
	mt, ok := v.(map[string]interface{})
	if rt == "" && ok && mt[attrCtx] != nil && mt[attrCtx].(string) == "" {
		return nil, ErrUnknownType
	}

	m := map[string]interface{}{}
	if err := Encode(v, &m, true); err != nil {
		return nil, err
	}

	tm, err := TypeMap(m)
	if err != nil {
		return nil, err
	}

	b := []byte{}
	enc := codec.NewEncoderBytes(&b, CborHandler())
	if err := enc.Encode(tm); err != nil {
		return nil, err
	}

	return b, nil
}
