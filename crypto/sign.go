package crypto

import "nimona.io/go/encoding"

// Sign block (container) with given key and return a signature block (container)
func Sign(v interface{}, key *Key) (*Signature, error) {
	b, err := encoding.Marshal(v)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	if err := encoding.UnmarshalSimple(b, &m); err != nil {
		return nil, err
	}

	// TODO replace ES256 with OH that should deal with removing the @sig
	delete(m, "@sig")

	b, err = encoding.Marshal(m)
	if err != nil {
		return nil, err
	}

	return NewSignature(key, "ES256", b)
}
