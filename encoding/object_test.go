package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testTypeFoo       = "test/foo"
	testTypeComposite = "test/composite"
)

type TestFoo struct {
	Foo string `json:"foo"`
}

type TestComposite struct {
	Foo       string  `json:"foo"`
	Signer    *Object `json:"@signer"`
	RawObject *Object `json:"@"`
}

func TestObjectMap(t *testing.T) {
	Register(testTypeFoo, &TestFoo{})

	em := map[string]interface{}{
		"@ctx:s":          "test/message",
		"simple-string:s": "hello world",
		"nested-object:o": map[string]interface{}{
			"@ctx:s": "test/something-random",
			"foo:s":  "bar",
		},
		"@signer:O": map[string]interface{}{
			"@ctx:s": testTypeFoo,
			"crv:s":  "P-256",
			"kty:s":  "EC",
		},
	}

	o := NewObject(em)
	assert.NotNil(t, o)

	assert.NotNil(t, o.SignerKey())
	assert.IsType(t, o.SignerKey(), &Object{})

	m := o.Map()
	assert.Equal(t, em, m)
}

func TestCompositeObjectMap(t *testing.T) {
	Register(testTypeFoo, &TestFoo{})
	Register(testTypeComposite, &TestComposite{})

	em := map[string]interface{}{
		"@ctx:s": testTypeComposite,
		"foo:s":  "hello world",
		"@signer:O": map[string]interface{}{
			"@ctx:s": testTypeFoo,
			"crv:s":  "P-256",
			"kty:s":  "EC",
		},
	}

	o := NewObject(em)
	assert.NotNil(t, o)

	assert.NotNil(t, o.SignerKey())
	assert.IsType(t, &Object{}, o.SignerKey())

	m := o.Map()
	assert.Equal(t, em, m)
}
