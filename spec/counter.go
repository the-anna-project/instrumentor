package spec

// Counter is a metric that can be arbitrarily incremented.
type Counter interface {
	// Increment increments the current counter by the given delta.
	Increment(delta float64) error
	IncrementWithLabels(delta float64, values ...string) error
}

type CounterConfig interface {
	Help() string
	Labels() []string
	Name() string
	SetHelp(string)
	SetLabels([]string)
	SetName(string)
}
