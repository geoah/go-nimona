// Code generated by nimona.io/go/cmd/objectify. DO NOT EDIT.

// +build !generate

package crypto

import (
	"nimona.io/go/encoding"
)

// ToMap returns a map compatible with f12n
func (s Key) ToMap() map[string]interface{} {

	m := map[string]interface{}{
		"@ctx:s":    "/key",
		"alg:s":     s.Algorithm,
		"kid:s":     s.KeyID,
		"kty:s":     s.KeyType,
		"use:s":     s.KeyUsage,
		"key_ops:s": s.KeyOps,
		"x5c:s":     s.X509CertChain,
		"x5t:s":     s.X509CertThumbprint,
		"x5tS256:s": s.X509CertThumbprintS256,
		"x5u:s":     s.X509URL,
		"crv:s":     s.Curve,
		"x:d":       s.X,
		"y:d":       s.Y,
		"d:d":       s.D,
	}
	return m
}

// ToObject returns a f12n object
func (s Key) ToObject() *encoding.Object {
	return encoding.NewObjectFromMap(s.ToMap())
}

// FromMap populates the struct from a f12n compatible map
func (s *Key) FromMap(m map[string]interface{}) error {
	if v, ok := m["alg:s"].(string); ok {
		s.Algorithm = v
	}
	if v, ok := m["kid:s"].(string); ok {
		s.KeyID = v
	}
	if v, ok := m["kty:s"].(string); ok {
		s.KeyType = v
	}
	if v, ok := m["use:s"].(string); ok {
		s.KeyUsage = v
	}
	if v, ok := m["key_ops:s"].(string); ok {
		s.KeyOps = v
	}
	if v, ok := m["x5c:s"].(string); ok {
		s.X509CertChain = v
	}
	if v, ok := m["x5t:s"].(string); ok {
		s.X509CertThumbprint = v
	}
	if v, ok := m["x5tS256:s"].(string); ok {
		s.X509CertThumbprintS256 = v
	}
	if v, ok := m["x5u:s"].(string); ok {
		s.X509URL = v
	}
	if v, ok := m["crv:s"].(string); ok {
		s.Curve = v
	}
	if v, ok := m["x:d"].([]byte); ok {
		s.X = v
	}
	if v, ok := m["y:d"].([]byte); ok {
		s.Y = v
	}
	if v, ok := m["d:d"].([]byte); ok {
		s.D = v
	}
	s.RawObject = encoding.NewObjectFromMap(m)
	return nil
}

// FromObject populates the struct from a f12n object
func (s *Key) FromObject(o *encoding.Object) error {
	return s.FromMap(o.Map())
}

// GetType returns the object's type
func (s Key) GetType() string {
	return "/key"
}
