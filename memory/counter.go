package memory

import (
	objectspec "github.com/the-anna-project/spec/object"
)

// CounterConfig represents the configuration used to create a new memory
// counter object.
type CounterConfig struct {
}

// DefaultCounterConfig provides a default configuration to create a new memory
// counter object by best effort.
func DefaultCounterConfig() CounterConfig {
	newConfig := CounterConfig{}

	return newConfig
}

// NewCounter creates a new configured memory counter object.
func NewCounter(config CounterConfig) (objectspec.InstrumentorCounter, error) {
	newCounter := &counter{
		CounterConfig: config,
	}

	return newCounter, nil
}

type counter struct {
	CounterConfig
}

func (c *counter) IncrBy(delta float64) {
}
