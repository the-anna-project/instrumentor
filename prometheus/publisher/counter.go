package publisher

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/the-anna-project/instrumentor/spec"
)

// CounterConfig represents the configuration used to create a new prometheus
// publisher counter object.
type CounterConfig struct {
	// Settings.

	// help represents some sort of informative description of the registered
	// metric.
	help string
	// labels represents labels as being used by prometheus. They partition
	// metrics.
	labels []string
	// name represents the metric's key as it is supposed to be registered. In the
	// scope of prometheus publisher  this is expected to be an underscored
	// string.
	name string
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

// DefaultCounterConfig provides a default configuration to create a new
// prometheus publisher counter object by best effort.
func DefaultCounterConfig() *CounterConfig {
	newConfig := &CounterConfig{
		// Settings.
		help:   "",
		labels: nil,
		name:   "",
	}

	return newConfig
}

// NewCounter creates a new configured prometheus publisher counter object.
func NewCounter(config spec.CounterConfig) (*Counter, error) {
	// Settings.
	if config.Help() == "" {
		return nil, maskAnyf(invalidConfigError, "help must not be empty")
	}
	if config.Name() == "" {
		return nil, maskAnyf(invalidConfigError, "name must not be empty")
	}

	var clientCounter prometheus.Counter
	var clientCounterVec *prometheus.CounterVec

	if len(config.Labels()) == 0 {
		clientCounter = prometheus.NewCounter(
			prometheus.CounterOpts{
				Help: config.Help(),
				Name: config.Name(),
			},
		)
	} else {
		clientCounterVec = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Help: config.Help(),
				Name: config.Name(),
			},
			config.Labels(),
		)
	}

	newCounter := &Counter{
		ClientCounter:    clientCounter,
		ClientCounterVec: clientCounterVec,
	}

	return newCounter, nil
}

type Counter struct {
	// Public.
	ClientCounter    prometheus.Counter
	ClientCounterVec *prometheus.CounterVec
}

func (c *Counter) Increment(delta float64) error {
	if c.ClientCounter == nil {
		// This error indicates that the counter has been configured with labels.
		// Therefore Counter.ObserveWithLabels must be used.
		return maskAnyf(invalidConfigError, "counter must be configured")
	}

	c.ClientCounter.Add(delta)

	return nil
}

func (c *Counter) IncrementWithLabels(delta float64, values ...string) error {
	if c.ClientCounterVec == nil {
		// This error indicates that the counter has not been configured with
		// labels. Therefore Counter.Observe must be used.
		return maskAnyf(invalidConfigError, "counter must be configured")
	}
	if len(values) == 0 {
		return maskAnyf(invalidConfigError, "labels must not be empty")
	}

	c.ClientCounterVec.WithLabelValues(values...).Add(delta)

	return nil
}
