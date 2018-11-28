package encoding

import (
	"nimona.io/go/base58"
)

// Object for everything f12n
type Object struct {
	data      map[string]interface{}
	ctx       string
	policy    *Object
	authority *Object
	signer    *Object
	signature *Object
}

// NewObjectFromStruct returns an object from a struct
func NewObjectFromStruct(v interface{}) (*Object, error) {
	m, err := Encode(v, true)
	if err != nil {
		return nil, err
	}

	tm, err := TypeMap(m)
	if err != nil {
		return nil, err
	}

	return NewObject(tm), nil
}

// NewObject returns an object from a map
func NewObject(m map[string]interface{}) *Object {
	o := &Object{
		data: map[string]interface{}{},
	}
	for k, v := range m {
		o.SetRaw(k, v)
	}
	return o
}

// Materialize returns a populated struct based on the registered type
// func (o *Object) Materialize() (interface{}, error) {
// 	m := o.Map()
// 	um, err := UntypeMap(m)
// 	if err != nil {
// 		return nil, err
// 	}

// 	v := GetInstance(GetType(um))
// 	if v == nil {
// 		v = &map[string]interface{}{}
// 	}
// 	if err := Decode(um, v, true); err != nil {
// 		return nil, err
// 	}

// 	return v, nil
// }

// Hash returns the object's hash
func (o *Object) Hash() []byte {
	return Hash(o)
}

// HashBase58 returns the object's hash base58 encoded
func (o *Object) HashBase58() string {
	return base58.Encode(Hash(o))
}

// Map returns the object as a map
func (o *Object) Map() map[string]interface{} {
	return o.data
}

// Type returns the object's type
func (o *Object) Type() string {
	return o.ctx
}

// SetType sets the object's type
func (o *Object) SetType(v string) {
	o.ctx = v
}

// Signature returns the object's signature, or nil
func (o *Object) Signature() *Object {
	return o.signature
}

// SetSignature sets the object's signature
func (o *Object) SetSignature(v *Object) {
	o.signature = v
}

// AuthorityKey returns the object's creator, or nil
func (o *Object) AuthorityKey() *Object {
	return o.authority
}

// SetAuthorityKey sets the object's creator
func (o *Object) SetAuthorityKey(v *Object) {
	o.authority = v
}

// SignerKey returns the object's signer, or nil
func (o *Object) SignerKey() *Object {
	return o.signer
}

// SetSignerKey sets the object's signer
func (o *Object) SetSignerKey(v *Object) {
	o.signer = v
}

// Policy returns the object's policy, or nil
func (o *Object) Policy() *Object {
	return o.policy
}

// SetPolicy sets the object's policy
func (o *Object) SetPolicy(v *Object) {
	o.policy = v
}

// GetRaw -
func (o *Object) GetRaw(lk string) interface{} {
	// TODO(geoah) do we need to verify type if k has hint?
	lk = getCleanKeyName(lk)
	for k, v := range o.data {
		if getCleanKeyName(k) == lk {
			return v
		}
	}

	return nil
}

// SetRaw -
func (o *Object) SetRaw(k string, v interface{}) {
	// add type hint if not already set
	et := getFullType(k)
	if et == "" {
		k += ":" + getTypeHint(v)
	}

	// clear the signature as it has been invalidated
	delete(o.data, "@sig")
	o.signature = nil

	// add the attribute in the data map
	o.data[k] = v

	// check if this is a magic attribute and set it
	ck := getCleanKeyName(k)
	switch ck {
	case "@ctx":
		t, ok := v.(string)
		if !ok {
			panic("invalid type for @ctx")
		}
		o.ctx = t
	case "@policy":
		m, ok := v.(map[string]interface{})
		if !ok {
			panic("invalid type for @policy")
		}
		o.policy = NewObject(m)
	case "@authority":
		m, ok := v.(map[string]interface{})
		if !ok {
			panic("invalid type for @authority")
		}
		o.authority = NewObject(m)
	case "@signer":
		m, ok := v.(map[string]interface{})
		if !ok {
			panic("invalid type for @signer")
		}
		o.signer = NewObject(m)
	case "@sig":
		m, ok := v.(map[string]interface{})
		if !ok {
			panic("invalid type for @sig")
		}
		o.signature = NewObject(m)
	}
}
