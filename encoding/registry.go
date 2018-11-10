package encoding

import (
	"reflect"
	"sync"
)

var (
	registry = sync.Map{}
)

// Register registers structs we can decode
func Register(t string, v interface{}) {
	registry.Store(t, reflect.TypeOf(v))
}

// GetInstance from string type
func GetInstance(t string) interface{} {
	rt, ok := registry.Load(t)
	if !ok {
		return nil
	}
	return reflect.New(rt.(reflect.Type).Elem()).Interface()
}

// GetType from interface
func GetType(t reflect.Type) string {
	var rt string
	registry.Range(func(k, v interface{}) bool {
		if v.(reflect.Type) == t {
			rt = k.(string)
		}
		return true
	})
	return rt
}
