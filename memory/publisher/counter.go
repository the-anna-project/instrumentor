package publisher

import (
	"github.com/the-anna-project/instrumentor/spec"
)

// CounterConfig represents the configuration used to create a new memory
// publisher counter object.
type CounterConfig struct {
	// Settings.
	help   string
	labels []string
	name   string
}

func (cc *CounterConfig) Help() string {
	return cc.help
}

func (cc *CounterConfig) Labels() []string {
	return cc.labels
}

func (cc *CounterConfig) Name() string {
	return cc.name
}

func (cc *CounterConfig) SetHelp(help string) {
	cc.help = help
}

func (cc *CounterConfig) SetLabels(labels []string) {
	cc.labels = labels
}

func (cc *CounterConfig) SetName(name string) {
	cc.name = name
}

// DefaultCounterConfig provides a default configuration to create a new memory
// publisher counter object by best effort.
func DefaultCounterConfig() *CounterConfig {
	return &CounterConfig{
		// Settings.
		help:   "",
		labels: nil,
		name:   "",
	}
}

// NewCounter creates a new configured prometheus publisher counter object.
func NewCounter(config spec.CounterConfig) (*Counter, error) {
	newCounter := &Counter{}

	return newCounter, nil
}

type Counter struct {
}

func (c *Counter) Increment(delta float64) error {
	return nil
}

func (c *Counter) IncrementWithLabels(delta float64, values ...string) error {
	return nil
}
