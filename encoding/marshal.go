package encoding

import (
	"errors"
	"reflect"

	"github.com/ugorji/go/codec"
)

var (
	ErrMissingType = errors.New("missing type")
)

func Marshal(v interface{}) ([]byte, error) {
	rt := GetType(reflect.TypeOf(v))
	mt, ok := v.(map[string]interface{})
	if rt == "" && ok && mt[attrCtx] != nil && mt[attrCtx].(string) == "" {
		return nil, ErrMissingType
	}

	m := map[string]interface{}{}
	if err := Encode(v, &m, true); err != nil {
		return nil, err
	}

	b := []byte{}
	enc := codec.NewEncoderBytes(&b, CborHandler())
	if err := enc.Encode(m); err != nil {
		return nil, err
	}

	return b, nil
}
