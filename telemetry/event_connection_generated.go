// Code generated by nimona.io/go/cmd/objectify. DO NOT EDIT.

// +build !generate

package telemetry

import (
	"nimona.io/go/encoding"
)

// ToMap returns a map compatible with f12n
func (s ConnectionEvent) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"@ctx:s":      "nimona.io/telemetry/connection",
		"direction:s": s.Direction,
	}
	return m
}

// ToObject returns a f12n object
func (s ConnectionEvent) ToObject() *encoding.Object {
	return encoding.NewObjectFromMap(s.ToMap())
}

// FromMap populates the struct from a f12n compatible map
func (s *ConnectionEvent) FromMap(m map[string]interface{}) error {
	if v, ok := m["direction:s"].(string); ok {
		s.Direction = v
	}
	return nil
}

// FromObject populates the struct from a f12n object
func (s *ConnectionEvent) FromObject(o *encoding.Object) error {
	return s.FromMap(o.Map())
}

// GetType returns the object's type
func (s ConnectionEvent) GetType() string {
	return "nimona.io/telemetry/connection"
}