package crypto

import (
	"errors"
)

// MandatePolicy for Mandate
type MandatePolicy struct {
	Description string   `json:"description"`
	Resources   []string `json:"resources"`
	Actions     []string `json:"actions"`
	Effect      string   `json:"effect"`
}

// Mandate to give authority to a aubject to perform certain actions on the
// authority's behalf
type Mandate struct {
	Authority *Key           `json:"authority"`
	Subject   *Key           `json:"subject"`
	Policy    *MandatePolicy `json:"policy"`
	Signature *Signature     `json:"@sig"`
}

// NewMandate returns a signed mandate given an authority key, a subject key,
// and a policy
func NewMandate(authority, subject *Key, policy *MandatePolicy) (*Mandate, error) {
	if authority == nil {
		return nil, errors.New("missing authority")
	}

	if subject == nil {
		return nil, errors.New("missing subject")
	}

	m := &Mandate{
		Authority: authority,
		Subject:   subject,
		Policy:    policy,
	}

	// o, err := encoding.NewObjectFromStruct(m)
	// if err != nil {
	// 	return nil, err
	// }

	// s, err := Sign(o, authority)
	// if err != nil {
	// 	return nil, err
	// }

	return m, nil
}
