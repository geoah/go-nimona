package encoding

import (
	"reflect"

	"github.com/ugorji/go/codec"
)

// CborHandler for un/marshaling blocks
func CborHandler() *codec.CborHandle {
	ch := &codec.CborHandle{}
	ch.Canonical = true
	ch.DeleteOnNilMapValue = true
	ch.MapType = reflect.TypeOf(map[string]interface{}(nil))
	ch.TimeRFC3339 = true
	return ch
}
