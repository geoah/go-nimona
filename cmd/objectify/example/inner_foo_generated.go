// Code generated by nimona.io/go/cmd/objectify. DO NOT EDIT.

// +build !generate

package example

import (
	"nimona.io/go/encoding"
)

// ToMap returns a map compatible with f12n
func (s InnerFoo) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"@ctx:s":      "test/inn",
		"inner_bar:s": s.InnerBar,
		"i:i":         s.I,
		"i8:i":        s.I8,
		"i16:i":       s.I16,
		"i32:i":       s.I32,
		"i64:i":       s.I64,
		"u:u":         s.U,
		"u8:u":        s.U8,
		"u16:u":       s.U16,
		"u32:u":       s.U32,
		"f32:f":       s.F32,
		"f64:f":       s.F64,
	}
	if s.MoreInnerFoos != nil {
		sMoreInnerFoos := []map[string]interface{}{}
		for _, v := range s.MoreInnerFoos {
			sMoreInnerFoos = append(sMoreInnerFoos, v.ToMap())
		}
		m["inner_foos:a<o>"] = sMoreInnerFoos
	}
	if s.Ai8 != nil {
		m["ai8:a<i>"] = s.Ai8
	}
	if s.Ai16 != nil {
		m["ai16:a<i>"] = s.Ai16
	}
	if s.Ai32 != nil {
		m["ai32:a<i>"] = s.Ai32
	}
	if s.Ai64 != nil {
		m["ai64:a<i>"] = s.Ai64
	}
	if s.Au16 != nil {
		m["au16:a<u>"] = s.Au16
	}
	if s.Au32 != nil {
		m["au32:a<u>"] = s.Au32
	}
	if s.Af32 != nil {
		m["af32:a<f>"] = s.Af32
	}
	if s.Af64 != nil {
		m["af64:a<f>"] = s.Af64
	}
	return m
}

// ToObject returns a f12n object
func (s InnerFoo) ToObject() *encoding.Object {
	return encoding.NewObjectFromMap(s.ToMap())
}

// FromMap populates the struct from a f12n compatible map
func (s *InnerFoo) FromMap(m map[string]interface{}) error {
	if v, ok := m["inner_bar:s"].(string); ok {
		s.InnerBar = v
	}
	s.MoreInnerFoos = []*InnerFoo{}
	if ss, ok := m["inner_foos:a<o>"].([]interface{}); ok {
		for _, si := range ss {
			if v, ok := si.(map[string]interface{}); ok {
				sMoreInnerFoos := &InnerFoo{}
				if err := sMoreInnerFoos.FromMap(v); err != nil {
					return err
				}
				s.MoreInnerFoos = append(s.MoreInnerFoos, sMoreInnerFoos)
			} else if v, ok := m["inner_foos:a<o>"].(*InnerFoo); ok {
				s.MoreInnerFoos = append(s.MoreInnerFoos, v)
			}
		}
	}
	if v, ok := m["i:i"].(int); ok {
		s.I = v
	}
	if v, ok := m["i8:i"].(int8); ok {
		s.I8 = v
	}
	if v, ok := m["i16:i"].(int16); ok {
		s.I16 = v
	}
	if v, ok := m["i32:i"].(int32); ok {
		s.I32 = v
	}
	if v, ok := m["i64:i"].(int64); ok {
		s.I64 = v
	}
	if v, ok := m["u:u"].(uint); ok {
		s.U = v
	}
	if v, ok := m["u8:u"].(uint8); ok {
		s.U8 = v
	}
	if v, ok := m["u16:u"].(uint16); ok {
		s.U16 = v
	}
	if v, ok := m["u32:u"].(uint32); ok {
		s.U32 = v
	}
	if v, ok := m["f32:f"].(float32); ok {
		s.F32 = v
	}
	if v, ok := m["f64:f"].(float64); ok {
		s.F64 = v
	}
	if v, ok := m["ai8:a<i>"].([]int8); ok {
		s.Ai8 = v
	}
	if v, ok := m["ai16:a<i>"].([]int16); ok {
		s.Ai16 = v
	}
	if v, ok := m["ai32:a<i>"].([]int32); ok {
		s.Ai32 = v
	}
	if v, ok := m["ai64:a<i>"].([]int64); ok {
		s.Ai64 = v
	}
	if v, ok := m["au16:a<u>"].([]uint16); ok {
		s.Au16 = v
	}
	if v, ok := m["au32:a<u>"].([]uint32); ok {
		s.Au32 = v
	}
	if v, ok := m["af32:a<f>"].([]float32); ok {
		s.Af32 = v
	}
	if v, ok := m["af64:a<f>"].([]float64); ok {
		s.Af64 = v
	}
	return nil
}

// FromObject populates the struct from a f12n object
func (s *InnerFoo) FromObject(o *encoding.Object) error {
	return s.FromMap(o.Map())
}

// GetType returns the object's type
func (s InnerFoo) GetType() string {
	return "test/inn"
}
