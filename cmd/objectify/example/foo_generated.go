// Code generated by nimona.io/go/cmd/objectify. DO NOT EDIT.

// +build !generate

package example

import (
	"nimona.io/go/encoding"
)

// ToMap returns a map compatible with f12n
func (s Foo) ToMap() map[string]interface{} {
	sInnerFoos := []map[string]interface{}{}
	for _, v := range s.InnerFoos {
		sInnerFoos = append(sInnerFoos, v.ToMap())
	}

	m := map[string]interface{}{
		"@ctx:s":          "test/foo",
		"bar:s":           s.Bar,
		"bars:A<s>":       s.Bars,
		"inner_foo:O":     s.InnerFoo.ToMap(),
		"inner_foos:A<O>": sInnerFoos,
	}
	return m
}

// ToObject returns a f12n object
func (s Foo) ToObject() *encoding.Object {
	return encoding.NewObjectFromMap(s.ToMap())
}

// FromMap populates the struct from a f12n compatible map
func (s *Foo) FromMap(m map[string]interface{}) error {
	s.RawObject = encoding.NewObjectFromMap(m)
	if v, ok := m["bar:s"].(string); ok {
		s.Bar = v
	}
	if v, ok := m["bars:A<s>"].([]string); ok {
		s.Bars = v
	}
	if v, ok := m["inner_foo:O"].(map[string]interface{}); ok {
		s.InnerFoo = &InnerFoo{}
		if err := s.InnerFoo.FromMap(v); err != nil {
			return err
		}
	} else if v, ok := m["inner_foo:O"].(*InnerFoo); ok {
		s.InnerFoo = v
	}
	s.InnerFoos = []*InnerFoo{}
	if ss, ok := m["inner_foos:A<O>"].([]interface{}); ok {
		for _, si := range ss {
			if v, ok := si.(map[string]interface{}); ok {
				sInnerFoos := &InnerFoo{}
				if err := sInnerFoos.FromMap(v); err != nil {
					return err
				}
				s.InnerFoos = append(s.InnerFoos, sInnerFoos)
			} else if v, ok := m["inner_foos:A<O>"].(*InnerFoo); ok {
				s.InnerFoos = append(s.InnerFoos, v)
			}
		}
	}
	return nil
}

// FromObject populates the struct from a f12n object
func (s *Foo) FromObject(o *encoding.Object) error {
	return s.FromMap(o.Map())
}

// GetType returns the object's type
func (s Foo) GetType() string {
	return "test/foo"
}
