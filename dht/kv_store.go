package dht

import (
	"sync"

	"nimona.io/go/crypto"
	"nimona.io/go/encoding"
)

type Store struct {
	// TODO replace with async maps
	values    sync.Map
	providers sync.Map
	lock      sync.RWMutex
}

func newStore() (*Store, error) {
	s := &Store{
		values:    sync.Map{},
		providers: sync.Map{},
		lock:      sync.RWMutex{},
	}
	return s, nil
}

func (s *Store) PutProvider(provider *Provider) error {
	// TODO verify payload type
	b, _ := encoding.Marshal(provider)
	h := crypto.NewSha3(b).Base58()
	s.providers.Store(h, provider)
	return nil
}

func (s *Store) GetProviders(key string) ([]*Provider, error) {
	providers := []*Provider{}
	s.providers.Range(func(k, v interface{}) bool {
		provider := v.(*Provider)
		for _, id := range provider.BlockIDs {
			if id == key {
				providers = append(providers, provider)
				break
			}
		}
		return true
	})

	return providers, nil
}

// GetAllProviders returns all providers and the values they are providing
func (s *Store) GetAllProviders() ([]*Provider, error) {
	providers := []*Provider{}
	s.providers.Range(func(k, v interface{}) bool {
		providers = append(providers, v.(*Provider))
		return true
	})

	return providers, nil
}
