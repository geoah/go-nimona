package telemetry

// Collectable for metric events
type Collectable interface {
	Collection() string
	Measurements() map[string]interface{}
}

// ConnectionEvent for reporting connection info
//proteus:generate
type ConnectionEvent struct {
	// Event attributes
	Direction string `json:"direction"`
}

// Collection returns the string representation of the structure
func (ce *ConnectionEvent) Collection() string {
	return "nimona.io/telemetry.connection"
}

// Measurements returns a map with all the metrics for the event
func (ce *ConnectionEvent) Measurements() map[string]interface{} {
	return map[string]interface{}{
		"direction": ce.Direction,
	}
}

// BlockEvent for reporting block metrics
//proteus:generate
type BlockEvent struct {
	// Event attributes
	Direction   string `json:"direction"`
	ContentType string `json:"contentType"`
	BlockSize   int    `json:"size"`
	// Signature   *crypto.Signature `json:"-"`
}

// Collection returns the string representation of the structure
func (ee *BlockEvent) Collection() string {
	return "nimona.io/telemetry.block"
}

// Measurements returns a map with all the metrics for the event
func (ee *BlockEvent) Measurements() map[string]interface{} {
	return map[string]interface{}{
		"direction":    ee.Direction,
		"content_type": ee.ContentType,
		"block_size":   ee.BlockSize,
	}
}
