package spec

// Histogram is a metric to observe samples over time.
type Histogram interface {
	// Observe tracks the given sample used for aggregation of the current
	// histogramm.
	Observe(sample float64) error
	ObserveWithLabels(sample float64, values ...string) error
}

type HistogramConfig interface {
	Buckets() []float64
	Help() string
	Labels() []string
	Name() string
	SetBuckets([]float64)
	SetHelp(string)
	SetLabels([]string)
	SetName(string)
}
