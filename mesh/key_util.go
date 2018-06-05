package mesh

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/nacl/box"
)

// LoadOrCreateLocalPeerInfo from/to a JSON encoded file
func (reg *registry) LoadOrCreateLocalPeerInfo(path string) (*SecretPeerInfo, error) {
	if path == "" {
		return nil, errors.New("missing key path")
	}

	if _, err := os.Stat(path); err == nil {
		return reg.LoadSecretPeerInfo(path)
	}

	log.Printf("* Key path does not exist, creating new key in '%s'\n", path)

	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	pi := &SecretPeerInfo{
		PeerInfo: PeerInfo{
			Addresses: []string{},
			PublicKey: *pub,
		},
		SecretKey: *priv,
	}

	pi.ID = fmt.Sprintf("%x", pi.GetPublicKey().ToKID())

	reg.keyring.ImportBoxKey(pub, priv)

	if err := reg.StoreSecretPeerInfo(pi, path); err != nil {
		return nil, err
	}

	return pi, nil
}

// CreateNewPeer with a new generated key, mostly used for testing
func (reg *registry) CreateNewPeer() (*SecretPeerInfo, error) {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	reg.keyring.ImportBoxKey(pub, priv)

	spi := &SecretPeerInfo{
		PeerInfo: PeerInfo{
			Addresses: []string{},
			PublicKey: *pub,
		},
		SecretKey: *priv,
	}

	spi.ID = fmt.Sprintf("%x", spi.GetPublicKey().ToKID())

	return spi, nil
}

// LoadSecretPeerInfo from a JSON encoded file
func (reg *registry) LoadSecretPeerInfo(path string) (*SecretPeerInfo, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pi := &SecretPeerInfo{}
	if err := json.Unmarshal(raw, &pi); err != nil {
		return nil, err
	}

	reg.keyring.ImportBoxKey(&pi.PublicKey, &pi.SecretKey)

	return pi, nil
}

// StoreSecretPeerInfo to a JSON encoded file
func (reg *registry) StoreSecretPeerInfo(pi *SecretPeerInfo, path string) error {
	raw, err := json.Marshal(pi)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, raw, 0644)
}
