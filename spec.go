package instrumentor

import (
	"net/http"
)

// Counter is a metric that can be arbitrarily incremented.
type Counter interface {
	// IncrBy increments the current counter by the given delta.
	IncrBy(delta float64)
}

// Gauge is a metric that can be arbitrarily incremented or
// decremented.
type Gauge interface {
	// DecrBy decrements the current gauge by the given delta.
	DecrBy(delta float64)
	// IncrBy increments the current gauge by the given delta.
	IncrBy(delta float64)
}

// Histogram is a metric to observe samples over time.
type Histogram interface {
	// Observe tracks the given sample used for aggregation of the current
	// histogramm.
	Observe(sample float64)
}

// Service represents an abstraction of instrumentation libraries to manage
// application metrics.
type Service interface {
	// ExecFunc wraps basic instrumentation around the given action and executes
	// it.
	//
	// The wrapped action causes the following metric's to be emitted. <prefix>
	// is described by the configured prefix of the current instrumentor.
	//
	//     <prefix>_<key>_durations_histogram_milliseconds
	//
	//         Holds the action's duration in milliseconds. This metric is
	//         emitted for each executed action.
	//
	//     <prefix>_<key>_errors_total
	//
	//         Holds the action's error count. This metric is emitted for each
	//         error returned by the given action.
	//
	ExecFunc(key string, action func() error) error
	// Counter provides a Counter for the given key. In case there does no counter
	// exist for the given key, one is created.
	Counter(key string) (Counter, error)
	// Gauge provides a Gauge for the given key. In case there does no counter
	// exist for the given key, one is created.
	Gauge(key string) (Gauge, error)
	// Gauge provides a Gauge for the given key. In case there does no counter
	// exist for the given key, one is created.
	Histogram(key string) (Histogram, error)
	// HTTPEndpoint returns the instrumentor's metric endpoint supposed to be
	// registered in the HTTP server using the instrumentor's HTTP handler.
	HTTPEndpoint() string
	// HTTPHandler returns the instrumentor's HTTP handler supposed to be
	// registered in the HTTP server using the instrumentor's HTTP endpoint.
	HTTPHandler() http.Handler
	// Prefixes returns the instrumentor's configured prefix.
	Prefixes() []string
	// NewKey returns a new metrics key having all configured prefixes and all
	// given parts properly joined. This could happen e.g. using underscores. In
	// this case it would look as follows.
	//
	//     prefix_prefix_s_s_s_s
	//
	NewKey(s ...string) string
	// WrapFunc wraps basic instrumentation around the given action. The returned
	// function can be used as e.g. retry action.
	//
	// The wrapped action causes the same metric's to be emitted as WrapFunc.
	WrapFunc(key string, action func() error) func() error
}
