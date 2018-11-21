package encoding

import (
	"errors"
	"fmt"
	"reflect"
)

// UntypeMap checks the type hints are correct, and removes them from the keys
func UntypeMap(m map[string]interface{}) (map[string]interface{}, error) {
	out := map[string]interface{}{}
	for k, v := range m {
		t := reflect.TypeOf(v)
		h := getTypeHint(t)
		eh := getFullType(k)
		if eh == "" {
			// key doesn't have type
			// TODO(geoah) is this even allowed?
			out[k] = v
			continue
		}
		if h != eh {
			return nil, fmt.Errorf("type hinted as %s, but is %s", eh, h)
		}
		// TODO should we be using type checks here?
		switch h {
		case hintMap:
			m, ok := v.(map[string]interface{})
			if !ok {
				return nil, errors.New("untype only supports map[string]interface{} maps")
			}
			var err error
			v, err = UntypeMap(m)
			if err != nil {
				return nil, err
			}
		case hintArray + "<" + hintMap + ">":
			vs, ok := v.([]interface{})
			if !ok {
				return nil, errors.New("untype only supports []interface{} for A<O>")
			}
			ovs := []interface{}{}
			for _, v := range vs {
				m, ok := v.(map[string]interface{})
				if !ok {
					return nil, errors.New("untype only supports map[string]interface{} maps")
				}
				ov, err := UntypeMap(m)
				if err != nil {
					return nil, err
				}
				ovs = append(ovs, ov)
			}
			v = ovs
		}
		k = k[:len(k)-len(h)-1]
		out[k] = v
	}
	return out, nil
}
