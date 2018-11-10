package encoding

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

func getDecodeHook(addCtx bool) mapstructure.DecodeHookFuncType {
	return func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		// fmt.Println("------------")
		// fmt.Println("from", from, from.Kind())
		// fmt.Println("to", to, to.Kind())
		// fmt.Println("data",data)

		if !addCtx {
			// TODO(geoah) explain this hack
			addCtx = true
			return data, nil
		}

		// decoding map to struct
		if obj, ok := data.(map[string]interface{}); ok {
			if to.Kind() == reflect.Map {
				return data, nil
			}

			t, ok := obj[attrCtx].(string)
			if !ok {
				return data, nil
			}

			v := GetInstance(t)
			if v == nil {
				return data, nil
			}

			delete(obj, attrCtx)
			if err := Decode(data, v, true); err != nil {
				return nil, err
			}

			return v, nil
		}

		return data, nil
	}
}

// Decode is a wrapper for mapstructure's Decode with our decodeHook that allows
// decoding maps to structs using the value of our @ctx attribute
func Decode(from interface{}, to interface{}, addCtx bool) error {
	dc := &mapstructure.DecoderConfig{
		Metadata:         &mapstructure.Metadata{},
		DecodeHook:       getDecodeHook(addCtx),
		Result:           to,
		TagName:          "json",
		WeaklyTypedInput: true,
	}
	dec, err := mapstructure.NewDecoder(dc)
	if err != nil {
		return err
	}

	return dec.Decode(from)
}
