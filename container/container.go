package container

import "nimona.io/go/crypto"

// Container wraps all objects
type Container map[string]interface{}

// Signature returns the container's signature, or nil
func (c Container) Signature() *crypto.Signature {
	return c["@"].(*crypto.Signature)
}

// Creator returns the container's creator, or nil
func (c Container) Creator() *crypto.Key {
	return c["@id.key"].(*crypto.Key) // TODO(geoah) return id's pk
}

// Signer returns the container's signer, or nil
func (c Container) Signer() *crypto.Key {
	return c["@id.peer.key"].(*crypto.Key) // TODO(geoah) return peer's id
}
