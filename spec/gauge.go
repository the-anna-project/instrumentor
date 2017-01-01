package spec

// Gauge is a metric that can be arbitrarily incremented or
// decremented.
type Gauge interface {
	// Decrement decrements the current gauge by the given delta.
	Decrement(delta float64) error
	DecrementWithLabels(delta float64, values ...string) error
	// Increment increments the current gauge by the given delta.
	Increment(delta float64) error
	IncrementWithLabels(delta float64, values ...string) error
	Set(value float64) error
	SetWithLabels(value float64, values ...string) error
}

type GaugeConfig interface {
	Help() string
	Labels() []string
	Name() string
	SetHelp(string)
	SetLabels([]string)
	SetName(string)
}
