package crypto

import "nimona.io/go/encoding"

func Sign(v interface{}, key *Key) (*Signature, error) {
	b, err := encoding.Marshal(v)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	if err := encoding.UnmarshalSimple(b, &m); err != nil {
		return nil, err
	}

	delete(m, "@sig")
	b, err = encoding.Marshal(m)
	if err != nil {
		return nil, err
	}

	return NewSignature(key, ES256, b)
}
