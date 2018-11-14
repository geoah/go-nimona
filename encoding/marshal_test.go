package encoding

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"nimona.io/go/base58"
)

const (
	testTypeKey       = "type:key"
	testTypeSignature = "type:sig"
	testTypeMessage   = "type:msg"
)

type TestKey struct {
	Algorithm string `json:"alg,omitempty"`
	X         []byte `json:"x,omitempty"`
	Y         []byte `json:"y,omitempty"`
	D         []byte `json:"d,omitempty"`
}

type TestSignature struct {
	Key       *TestKey `json:"key"`
	Alg       string   `json:"alg"`
	Signature []byte   `json:"sig"`
}

type TestMessage struct {
	Body      string         `json:"body"`
	Timestamp string         `json:"timestamp"`
	Signature *TestSignature `json:"@sig"`
}

func TestMarshalUnmarshal(t *testing.T) {
	Register(testTypeKey, &TestKey{})
	Register(testTypeSignature, &TestSignature{})
	Register(testTypeMessage, &TestMessage{})

	ek := &TestKey{
		Algorithm: "a",
		X:         []byte{1, 2, 3},
		Y:         []byte{4, 5, 6},
		D:         []byte{7, 8, 9},
	}

	es := &TestSignature{
		Key:       ek,
		Alg:       "b",
		Signature: []byte{10, 11, 12},
	}

	em := &TestMessage{
		Body:      "hello",
		Timestamp: "2018-11-09T22:07:21Z",
		Signature: es,
	}

	bs, err := Marshal(em)
	assert.NoError(t, err)

	assert.Equal(t, "dyWcbruKbSa1h4ZWvBQQPuPoJmfXw4EqQePCqVXDJe264bmUyAhMg1wTS"+
		"aFtMZsE9Zy24Hcn5NJyGpoouvAZWrkno5Wnud97mSDHn55zqg5oRxFt48qyAS3FzpGbjB"+
		"dmVeFkauQRs79JaF4FvC4SGFQdTm8SrcH5UUtXSPcZuB5CU6ifTktWqa9", base58.Encode(bs))

	m := &TestMessage{}
	err = UnmarshalInto(bs, m)
	assert.NoError(t, err)
	assert.Equal(t, em, m)

	v, err := Unmarshal(bs)
	assert.NoError(t, err)
	assert.Equal(t, em, v)

	xxx, _ := json.Marshal(m)
	fmt.Println(string(xxx))
}

func TestMarshalUnmarshalToMap(t *testing.T) {
	Register(testTypeKey, &TestKey{})
	Register(testTypeSignature, &TestSignature{})
	Register(testTypeMessage, &TestMessage{})

	ek := &TestKey{
		Algorithm: "a",
		X:         []byte{1, 2, 3},
		Y:         []byte{4, 5, 6},
		D:         []byte{7, 8, 9},
	}

	es := struct {
		Type      string   `json:"@ctx"`
		Key       *TestKey `json:"key"`
		Alg       string   `json:"alg"`
		Signature []byte   `json:"sig"`
	}{
		Type:      "type:other",
		Key:       ek,
		Alg:       "b",
		Signature: []byte{10, 11, 12},
	}

	bs, err := Marshal(es)
	assert.NoError(t, err)

	fmt.Println(base58.Encode(bs))

	em := map[string]interface{}{
		"@ctx": "type:other",
		"sig":  []byte{10, 11, 12},
		"alg":  "b",
		"key":  ek,
	}

	m := map[string]interface{}{}
	err = UnmarshalInto(bs, &m)
	assert.NoError(t, err)

	assert.Equal(t, em, m)

	xxx, _ := json.MarshalIndent(m, "", "")
	fmt.Println(string(xxx))
}
