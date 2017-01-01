package publisher

import (
	"github.com/the-anna-project/instrumentor/spec"
)

// HistogramConfig represents the configuration used to create a new memory
// publisher histogram.
type HistogramConfig struct {
	// Settings.
	buckets []float64
	help    string
	labels  []string
	name    string
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
// memory publisher histogram by best effort.
func DefaultHistogramConfig() *HistogramConfig {
	return &HistogramConfig{
		// Settings.
		buckets: nil,
		help:    "",
		labels:  nil,
		name:    "",
	}
}

// NewHistogram creates a new configured memory publisher histogram.
func NewHistogram(config spec.HistogramConfig) (*Histogram, error) {
	newHistogram := &Histogram{}

	return newHistogram, nil
}

type Histogram struct {
}

func (h *Histogram) Observe(sample float64) error {
	return nil
}

func (h *Histogram) ObserveWithLabels(sample float64, values ...string) error {
	return nil
}
