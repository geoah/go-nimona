package net

import (
	"nimona.io/go/encoding"
)

func init() {
	encoding.Register("/block.response", &BlockRequest{})
	encoding.Register("/block.response", &BlockResponse{})

	encoding.Register("/block.forward.request", &ForwardRequest{})
}
