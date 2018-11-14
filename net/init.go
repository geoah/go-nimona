package net

import (
	"nimona.io/go/encoding"
)

func init() {
	encoding.Register("/block.response", &BlockRequest{})
	encoding.Register("/block.response", &BlockResponse{})

	encoding.Register("/block.forward.request", &ForwardRequest{})

	encoding.Register("/handshake.syn", &HandshakeSyn{})
	encoding.Register("/handshake.syn-ack", &HandshakeSynAck{})
	encoding.Register("/handshake.ack", &HandshakeAck{})
}
