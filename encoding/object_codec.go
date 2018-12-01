package encoding

import (
	"github.com/davecgh/go-spew/spew"
	ucodec "github.com/ugorji/go/codec"
)

// CodecDecodeSelf helper for cbor/json unmarshaling
func (o *Object) CodecDecodeSelf(dec *ucodec.Decoder) {
	m := map[string]interface{}{}
	dec.MustDecode(&m)
	spew.Dump(m)
	o.FromMap(m)
}

// CodecEncodeSelf helper for cbor/json marshaling
func (o *Object) CodecEncodeSelf(enc *ucodec.Encoder) {
	m := o.Map()
	enc.MustEncode(m)
}
