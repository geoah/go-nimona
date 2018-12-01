// Code generated by nimona.io/go/cmd/objectify. DO NOT EDIT.

// +build !generate

package dht

import (
	"nimona.io/go/crypto"
	"nimona.io/go/encoding"
	"nimona.io/go/peers"
)

// ToMap returns a map compatible with f12n
func (s ProviderResponse) ToMap() map[string]interface{} {
	sProviders := []map[string]interface{}{}
	for _, v := range s.Providers {
		sProviders = append(sProviders, v.ToMap())
	}
	sClosestPeers := []map[string]interface{}{}
	for _, v := range s.ClosestPeers {
		sClosestPeers = append(sClosestPeers, v.ToMap())
	}
	m := map[string]interface{}{
		"@ctx:s":      "nimona.io/dht/provider.response",
		"requestID:s": s.RequestID,
	}
	if s.Providers != nil {
		m["providers:A<O>"] = s.Providers.ToMap()
	}
	if s.ClosestPeers != nil {
		m["closestPeers:A<O>"] = s.ClosestPeers.ToMap()
	}
	if s.Signer != nil {
		m["@signer:O"] = s.Signer.ToMap()
	}
	if s.Authority != nil {
		m["@authority:O"] = s.Authority.ToMap()
	}
	if s.Signature != nil {
		m["@sig:O"] = s.Signature.ToMap()
	}
	return m
}

// ToObject returns a f12n object
func (s ProviderResponse) ToObject() *encoding.Object {
	return encoding.NewObjectFromMap(s.ToMap())
}

// FromMap populates the struct from a f12n compatible map
func (s *ProviderResponse) FromMap(m map[string]interface{}) error {
	if v, ok := m["requestID:s"].(string); ok {
		s.RequestID = v
	}
	s.Providers = []*Provider{}
	if ss, ok := m["providers:A<O>"].([]interface{}); ok {
		for _, si := range ss {
			if v, ok := si.(map[string]interface{}); ok {
				sProviders := &Provider{}
				if err := sProviders.FromMap(v); err != nil {
					return err
				}
				s.Providers = append(s.Providers, sProviders)
			} else if v, ok := m["providers:A<O>"].(*Provider); ok {
				s.Providers = append(s.Providers, v)
			}
		}
	}
	s.ClosestPeers = []*peers.PeerInfo{}
	if ss, ok := m["closestPeers:A<O>"].([]interface{}); ok {
		for _, si := range ss {
			if v, ok := si.(map[string]interface{}); ok {
				sClosestPeers := &peers.PeerInfo{}
				if err := sClosestPeers.FromMap(v); err != nil {
					return err
				}
				s.ClosestPeers = append(s.ClosestPeers, sClosestPeers)
			} else if v, ok := m["closestPeers:A<O>"].(*peers.PeerInfo); ok {
				s.ClosestPeers = append(s.ClosestPeers, v)
			}
		}
	}
	s.RawObject = encoding.NewObjectFromMap(m)
	if v, ok := m["@signer:O"].(map[string]interface{}); ok {
		s.Signer = &crypto.Key{}
		if err := s.Signer.FromMap(v); err != nil {
			return err
		}
	} else if v, ok := m["@signer:O"].(*crypto.Key); ok {
		s.Signer = v
	}
	if v, ok := m["@authority:O"].(map[string]interface{}); ok {
		s.Authority = &crypto.Key{}
		if err := s.Authority.FromMap(v); err != nil {
			return err
		}
	} else if v, ok := m["@authority:O"].(*crypto.Key); ok {
		s.Authority = v
	}
	if v, ok := m["@sig:O"].(map[string]interface{}); ok {
		s.Signature = &crypto.Signature{}
		if err := s.Signature.FromMap(v); err != nil {
			return err
		}
	} else if v, ok := m["@sig:O"].(*crypto.Signature); ok {
		s.Signature = v
	}
	return nil
}

// FromObject populates the struct from a f12n object
func (s *ProviderResponse) FromObject(o *encoding.Object) error {
	return s.FromMap(o.Map())
}

// GetType returns the object's type
func (s ProviderResponse) GetType() string {
	return "nimona.io/dht/provider.response"
}
