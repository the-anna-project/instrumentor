// Package prometheus implements
// github.com/the-anna-project/instrumentor.Service and provides instrumentation
// primitives to manage application metrics.
package prometheus

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/the-anna-project/instrumentor"
)

// Config represents the configuration used to create a new instrumentor
// service.
type Config struct {
	// Settings.
	HTTPEndpoint string
	HTTPHandler  http.Handler
	Prefixes     []string
}

// DefaultConfig provides a default configuration to create a new instrumentor
// service by best effort.
func DefaultConfig() Config {
	return Config{
		// Settings.
		HTTPEndpoint: "/metrics",
		HTTPHandler:  prometheus.Handler(),
		Prefixes:     []string{},
	}
}

// New creates a new instrumentor service.
func New(config Config) (instrumentor.Service, error) {
	// Settings.
	if config.HTTPEndpoint == "" {
		return nil, maskAnyf(invalidConfigError, "HTTP endpoint must not be empty")
	}
	if config.HTTPHandler == nil {
		return nil, maskAnyf(invalidConfigError, "HTTP handler must not be empty")
	}
	if config.Prefixes == nil {
		return nil, maskAnyf(invalidConfigError, "prefixes must not be empty")
	}

	newService := &service{
		// Internals.
		counters:   map[string]instrumentor.Counter{},
		gauges:     map[string]instrumentor.Gauge{},
		histograms: map[string]instrumentor.Histogram{},
		mutex:      sync.Mutex{},

		// Settings.
		httpEndpoint: config.HTTPEndpoint,
		httpHandler:  config.HTTPHandler,
		prefixes:     config.Prefixes,
	}

	return newService, nil
}

type service struct {
	// Internals.
	counters   map[string]instrumentor.Counter
	gauges     map[string]instrumentor.Gauge
	histograms map[string]instrumentor.Histogram
	mutex      sync.Mutex

	// Settings.

	// httpEndpoint represents the HTTP endpoint used to register the httpHandler.
	// In the context of Prometheus this is usually /metrics.
	httpEndpoint string
	// httpHandler represents the HTTP handler used to register the Prometheus
	// registry in the HTTP server.
	httpHandler http.Handler
	// prefixes represents the Instrumentor's ordered prefixes.
	prefixes []string
}

func (s *service) ExecFunc(key string, action func() error) error {
	h, err := s.Histogram(s.NewKey(key, "durations", "histogram", "milliseconds"))
	if err != nil {
		return maskAny(err)
	}
	c, err := s.Counter(s.NewKey(key, "errors", "counter", "total"))
	if err != nil {
		return maskAny(err)
	}

	start := time.Now()

	err = action()
	if err != nil {
		c.IncrBy(1)
		return maskAny(err)
	}

	stop := time.Now()
	sample := float64(stop.Sub(start).Nanoseconds() / 1000000)
	h.Observe(sample)

	return nil
}

func (s *service) Counter(key string) (instrumentor.Counter, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if c, ok := s.counters[key]; ok {
		return c, nil
	}

	newConfig := DefaultCounterConfig()
	newConfig.Name = key
	newConfig.Help = helpFor("Counter", key)
	newCounter, err := NewCounter(newConfig)
	if err != nil {
		return nil, maskAny(err)
	}

	err = prometheus.Register(newCounter.(*counter).ClientCounter)
	if err != nil {
		return nil, maskAny(err)
	}
	s.counters[key] = newCounter

	return newCounter, nil
}

func (s *service) Gauge(key string) (instrumentor.Gauge, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if g, ok := s.gauges[key]; ok {
		return g, nil
	}

	newConfig := DefaultGaugeConfig()
	newConfig.Name = key
	newConfig.Help = helpFor("Gauge", key)
	newGauge, err := NewGauge(newConfig)
	if err != nil {
		return nil, maskAny(err)
	}

	err = prometheus.Register(newGauge.(*gauge).ClientGauge)
	if err != nil {
		return nil, maskAny(err)
	}
	s.gauges[key] = newGauge

	return newGauge, nil
}

func (s *service) Histogram(key string) (instrumentor.Histogram, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if h, ok := s.histograms[key]; ok {
		return h, nil
	}

	newConfig := DefaultHistogramConfig()
	newConfig.Name = key
	newConfig.Help = helpFor("Histogram", key)
	newHistogram, err := NewHistogram(newConfig)
	if err != nil {
		return nil, maskAny(err)
	}

	err = prometheus.Register(newHistogram.(*histogram).ClientHistogram)
	if err != nil {
		return nil, maskAny(err)
	}
	s.histograms[key] = newHistogram

	return newHistogram, nil
}

func (s *service) HTTPEndpoint() string {
	return s.httpEndpoint
}

func (s *service) HTTPHandler() http.Handler {
	return s.httpHandler
}

func (s *service) Prefixes() []string {
	return s.prefixes
}

func (s *service) NewKey(str ...string) string {
	return strings.Join(append(s.prefixes, str...), "_")
}

func (s *service) WrapFunc(key string, action func() error) func() error {
	wrappedFunc := func() error {
		err := s.ExecFunc(key, action)
		if err != nil {
			return maskAny(err)
		}

		return nil
	}

	return wrappedFunc
}
