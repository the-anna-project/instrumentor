package spec

import (
	"net/http"
)

// Publisher represents a service to abstract instrumentation libraries to emit
// application metrics.
type Publisher interface {
	// Boot initializes and starts the whole service like booting a machine. The
	// call to Boot blocks until the service is completely initialized, so you
	// might want to call it in a separate goroutine.
	Boot()
	// Counter provides a Counter for the given key. In case there does no counter
	// exist for the given key, one is created.
	Counter(config CounterConfig) (Counter, error)
	CounterConfig() CounterConfig
	// Gauge provides a Gauge for the given key. In case there does no counter
	// exist for the given key, one is created.
	Gauge(config GaugeConfig) (Gauge, error)
	GaugeConfig() GaugeConfig
	// Gauge provides a Gauge for the given key. In case there does no counter
	// exist for the given key, one is created.
	Histogram(config HistogramConfig) (Histogram, error)
	HistogramConfig() HistogramConfig
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
	// The wrapped action causes the following metric's to be emitted. <prefix>
	// is described by the configured prefix of the current instrumentor.
	//
	//     <prefix>_<key>_milliseconds
	//
	//         Holds the action's duration in milliseconds. This metric is
	//         emitted for each executed action.
	//
	//     <prefix>_<key>_error_total
	//
	//         Holds the action's error count. This metric is emitted for each
	//         error returned by the given action.
	//
	WrapFunc(key string, action func() error) func() error
	// Shutdown ends all processes of the service like shutting down a machine.
	// The call to Shutdown blocks until the service is completely shut down, so
	// you might want to call it in a separate goroutine.
	Shutdown()
}
