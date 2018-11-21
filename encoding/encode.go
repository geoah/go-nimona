package encoding

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// Encode is a wrapper for mapstructure's Decode with our decodeHook that allows
// encoding structs to maps
// TODO move addCtx to an option
func Encode(from interface{}, to *map[string]interface{}, addCtx bool) error {
	dc := &mapstructure.DecoderConfig{
		Metadata:         &mapstructure.Metadata{},
		DecodeHook:       getEncodeHook(addCtx),
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

func getEncodeHook(addCtx bool) mapstructure.DecodeHookFuncType {
	return func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		// fmt.Println("------------")
		// fmt.Println("from", from, from.Kind())
		// fmt.Println("to", to, to.Kind())
		// fmt.Println("data", data)

		if !addCtx {
			// TODO(geoah) explain this hack
			addCtx = true
			return data, nil
		}

		// decoding registered struct -- forced to map
		// TODO(geoah) WTF This is insanely hacky -- isn't it?
		if GetType(from) != "" {
			// same as "decoding struct to map"
			if t := GetType(from); t != "" {
				m := map[string]interface{}{}
				if err := Encode(data, &m, false); err != nil {
					return nil, err
				}
				m[attrCtx] = t
				for k, v := range m {
					if vt := GetType(reflect.TypeOf(v)); vt != "" {
						vm := map[string]interface{}{}
						if err := Encode(v, &vm, true); err != nil {
							return nil, err
						}
						vm[attrCtx] = vt
						m[k] = vm
					}
				}
				return m, nil
			}
			return data, nil
		}

		// decoding unknown struct to map
		if to.Kind() == reflect.Map {
			m := map[string]interface{}{}
			if err := Encode(data, &m, false); err != nil {
				return nil, err
			}
			// fmt.Println("m", reflect.TypeOf(m), m)
			return m, nil
		}

		return data, nil
	}
}
