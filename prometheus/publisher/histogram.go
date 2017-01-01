package publisher

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/the-anna-project/instrumentor/spec"
)

// HistogramConfig represents the configuration used to create a new prometheus
// publisher histogram.
type HistogramConfig struct {
	// Settings.

	// buckets represents a list of time ranges in seconds. Observed samples are
	// put into their corresponding ranges.
	//
	// A bucket's unit MUST be second. The buckets list MUST be ordered
	// incrementally.
	//
	// The buckets need to be properly configured to match the use case of the
	// oberseved samples, otherwise the histogram becomes pretty useless. E.g.
	// mapping samples of 25 milliseconds into a 5 second bucket makes no sense.
	buckets []float64
	// help represents some sort of informative description of the registered
	// metric.
	help string
	// labels represents labels as being used by prometheus. They partition
	// metrics.
	labels []string
	// name represents the metric's key as it is supposed to be registered. In the
	// scope of prometheus publisher this is expected to be an underscored string.
	name string
}

func (hc *HistogramConfig) Buckets() []float64 {
	return hc.buckets
}

func (hc *HistogramConfig) Help() string {
	return hc.help
}

func (hc *HistogramConfig) Labels() []string {
	return hc.labels
}

func (hc *HistogramConfig) Name() string {
	return hc.name
}

func (hc *HistogramConfig) SetBuckets(buckets []float64) {
	hc.buckets = buckets
}

func (hc *HistogramConfig) SetHelp(help string) {
	hc.help = help
}

func (hc *HistogramConfig) SetLabels(labels []string) {
	hc.labels = labels
}

func (hc *HistogramConfig) SetName(name string) {
	hc.name = name
}

// DefaultHistogramConfig provides a default configuration to create a new
// prometheus publisher histogram by best effort.
func DefaultHistogramConfig() *HistogramConfig {
	return &HistogramConfig{
		// Settings.
		buckets: []float64{.001, .002, .003, .004, .005, .01, .02, .03, .04, .05, .1, .2, .3, .4, .5, 1, 2, 3, 4, 5, 10},
		help:    "",
		labels:  nil,
		name:    "",
	}
}

// NewHistogram creates a new configured prometheus publisher histogram.
func NewHistogram(config spec.HistogramConfig) (*Histogram, error) {
	// Settings.
	if config.Buckets() == nil {
		return nil, maskAnyf(invalidConfigError, "buckets must not be empty")
	}
	if len(config.Buckets()) < 1 {
		return nil, maskAnyf(invalidConfigError, "buckets must contain at least 1 value")
	}
	if config.Help() == "" {
		return nil, maskAnyf(invalidConfigError, "help must not be empty")
	}
	if config.Name() == "" {
		return nil, maskAnyf(invalidConfigError, "name must not be empty")
	}

	var clientHistogram prometheus.Histogram
	var clientHistogramVec *prometheus.HistogramVec

	if len(config.Labels()) == 0 {
		clientHistogram = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Buckets: config.Buckets(),
				Help:    config.Help(),
				Name:    config.Name(),
			},
		)
	} else {
		clientHistogramVec = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Buckets: config.Buckets(),
				Help:    config.Help(),
				Name:    config.Name(),
			},
			config.Labels(),
		)
	}

	newHistogram := &Histogram{
		ClientHistogram:    clientHistogram,
		ClientHistogramVec: clientHistogramVec,
	}

	return newHistogram, nil
}

type Histogram struct {
	// Public.
	ClientHistogram    prometheus.Histogram
	ClientHistogramVec *prometheus.HistogramVec
}

func (h *Histogram) Observe(sample float64) error {
	if h.ClientHistogram == nil {
		// This error indicates that the histogram has been configured with labels.
		// Therefore Histogram.ObserveWithLabels must be used.
		return maskAnyf(invalidConfigError, "histogram must be configured")
	}

	h.ClientHistogram.Observe(sample)

	return nil
}

func (h *Histogram) ObserveWithLabels(sample float64, values ...string) error {
	if h.ClientHistogramVec == nil {
		// This error indicates that the histogram has not been configured with
		// labels. Therefore Histogram.Observe must be used.
		return maskAnyf(invalidConfigError, "histogram must be configured")
	}
	if len(values) == 0 {
		return maskAnyf(invalidConfigError, "labels must not be empty")
	}

	h.ClientHistogramVec.WithLabelValues(values...).Observe(sample)

	return nil
}
