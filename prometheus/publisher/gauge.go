package publisher

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/the-anna-project/instrumentor/spec"
)

// GaugeConfig represents the configuration used to create a new prometheus
// publisher gauge.
type GaugeConfig struct {
	// Settings.

	// help represents some sort of informative description of the registered
	// metric.
	help string
	// labels represents labels as being used by prometheus. They partition
	// metrics.
	labels []string
	// name represents the metric's key as it is supposed to be registered. In the
	// scope of prometheus this is expected to be an underscored string.
	name string
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

// DefaultGaugeConfig provides a default configuration to create a new
// prometheus publisher gauge by best effort.
func DefaultGaugeConfig() *GaugeConfig {
	return &GaugeConfig{
		// Settings.
		help:   "",
		labels: nil,
		name:   "",
	}
}

// NewGauge creates a new configured prometheus publisher gauge.
func NewGauge(config spec.GaugeConfig) (*Gauge, error) {
	// Settings.
	if config.Help() == "" {
		return nil, maskAnyf(invalidConfigError, "help must not be empty")
	}
	if config.Name() == "" {
		return nil, maskAnyf(invalidConfigError, "name must not be empty")
	}

	var clientGauge prometheus.Gauge
	var clientGaugeVec *prometheus.GaugeVec

	if len(config.Labels()) == 0 {
		clientGauge = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Help: config.Help(),
				Name: config.Name(),
			},
		)
	} else {
		clientGaugeVec = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Help: config.Help(),
				Name: config.Name(),
			},
			config.Labels(),
		)
	}

	newGauge := &Gauge{
		ClientGauge:    clientGauge,
		ClientGaugeVec: clientGaugeVec,
	}

	return newGauge, nil
}

type Gauge struct {
	// Public.
	ClientGauge    prometheus.Gauge
	ClientGaugeVec *prometheus.GaugeVec
}

func (g *Gauge) Decrement(delta float64) error {
	if g.ClientGauge == nil {
		// This error indicates that the gauge has been configured with labels.
		// Therefore Gauge.ObserveWithLabels must be used.
		return maskAnyf(invalidConfigError, "gauge must be configured")
	}

	g.ClientGauge.Sub(delta)

	return nil
}

func (g *Gauge) DecrementWithLabels(delta float64, values ...string) error {
	if g.ClientGaugeVec == nil {
		// This error indicates that the gauge has not been configured with labels.
		// Therefore Gauge.Observe must be used.
		return maskAnyf(invalidConfigError, "gauge must be configured")
	}
	if len(values) == 0 {
		return maskAnyf(invalidConfigError, "labels must not be empty")
	}

	g.ClientGaugeVec.WithLabelValues(values...).Sub(delta)

	return nil
}

func (g *Gauge) Increment(delta float64) error {
	if g.ClientGauge == nil {
		// This error indicates that the gauge has been configured with labels.
		// Therefore Gauge.ObserveWithLabels must be used.
		return maskAnyf(invalidConfigError, "gauge must be configured")
	}

	g.ClientGauge.Add(delta)

	return nil
}

func (g *Gauge) IncrementWithLabels(delta float64, values ...string) error {
	if g.ClientGaugeVec == nil {
		// This error indicates that the gauge has not been configured with labels.
		// Therefore Gauge.Observe must be used.
		return maskAnyf(invalidConfigError, "gauge must be configured")
	}
	if len(values) == 0 {
		return maskAnyf(invalidConfigError, "labels must not be empty")
	}

	g.ClientGaugeVec.WithLabelValues(values...).Add(delta)

	return nil
}

func (g *Gauge) Set(value float64) error {
	if g.ClientGauge == nil {
		// This error indicates that the gauge has been configured with labels.
		// Therefore Gauge.ObserveWithLabels must be used.
		return maskAnyf(invalidConfigError, "gauge must be configured")
	}

	g.ClientGauge.Set(value)

	return nil
}

func (g *Gauge) SetWithLabels(value float64, values ...string) error {
	if g.ClientGaugeVec == nil {
		// This error indicates that the gauge has not been configured with labels.
		// Therefore Gauge.Observe must be used.
		return maskAnyf(invalidConfigError, "gauge must be configured")
	}
	if len(values) == 0 {
		return maskAnyf(invalidConfigError, "labels must not be empty")
	}

	g.ClientGaugeVec.WithLabelValues(values...).Set(value)

	return nil
}
