package telemetry

import (
	"nimona.io/go/encoding"
)

func init() {
	encoding.Register("/telemetry.connection", &ConnectionEvent{})
	encoding.Register("/telemetry.block", &BlockEvent{})
}
