package publisher

import (
	"github.com/the-anna-project/instrumentor/spec"
)

// GaugeConfig represents the configuration used to create a new memory
// publisher gauge.
type GaugeConfig struct {
	// Settings.
	help   string
	labels []string
	name   string
}

func (gc *GaugeConfig) Help() string {
	return gc.help
}

func (gc *GaugeConfig) Labels() []string {
	return gc.labels
}

func (gc *GaugeConfig) Name() string {
	return gc.name
}

func (gc *GaugeConfig) SetHelp(help string) {
	gc.help = help
}

func (gc *GaugeConfig) SetLabels(labels []string) {
	gc.labels = labels
}

func (gc *GaugeConfig) SetName(name string) {
	gc.name = name
}

// DefaultGaugeConfig provides a default configuration to create a new memory
// publisher gauge by best effort.
func DefaultGaugeConfig() *GaugeConfig {
	return &GaugeConfig{
		// Settings.
		help:   "",
		labels: nil,
		name:   "",
	}
}

// NewGauge creates a new configured memory publisher gauge.
func NewGauge(config spec.GaugeConfig) (*Gauge, error) {
	newGauge := &Gauge{}

	return newGauge, nil
}

type Gauge struct {
}

func (g *Gauge) Decrement(delta float64) error {
	return nil
}

func (g *Gauge) DecrementWithLabels(delta float64, values ...string) error {
	return nil
}

func (g *Gauge) Increment(delta float64) error {
	return nil
}

func (g *Gauge) IncrementWithLabels(delta float64, values ...string) error {
	return nil
}

func (g *Gauge) Set(value float64) error {
	return nil
}

func (g *Gauge) SetWithLabels(value float64, values ...string) error {
	return nil
}
