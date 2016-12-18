package memory

import (
	"github.com/the-anna-project/instrumentor"
)

// GaugeConfig represents the configuration used to create a new memory gauge
// object.
type GaugeConfig struct {
}

// DefaultGaugeConfig provides a default configuration to create a new memory
// gauge object by best effort.
func DefaultGaugeConfig() GaugeConfig {
	newConfig := GaugeConfig{}

	return newConfig
}

// NewGauge creates a new configured memory gauge object.
func NewGauge(config GaugeConfig) (instrumentor.Gauge, error) {
	newGauge := &gauge{
		GaugeConfig: config,
	}

	return newGauge, nil
}

type gauge struct {
	GaugeConfig
}

func (g *gauge) DecrBy(delta float64) {
}

func (g *gauge) IncrBy(delta float64) {
}
