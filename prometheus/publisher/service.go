// Package publisher implements
// github.com/the-anna-project/instrumentor.Publisher and provides
// instrumentation primitives to emit application metrics.
package publisher

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/the-anna-project/instrumentor/spec"
)

// ServiceConfig represents the configuration used to create a new prometheus
// publisher service.
type ServiceConfig struct {
	// Settings.
	HTTPEndpoint string
	HTTPHandler  http.Handler
	Prefixes     []string
}

// DefaultServiceConfig provides a default configuration to create a new
// prometheus publisher service by best effort.
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		// Settings.
		HTTPEndpoint: "/metrics",
		HTTPHandler:  prometheus.Handler(),
		Prefixes:     []string{},
	}
}

// NewService creates a new prometheus publisher service.
func NewService(config ServiceConfig) (*Service, error) {
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

	newService := &Service{
		// Internals.
		closer:       make(chan struct{}, 1),
		counters:     map[string]*Counter{},
		bootOnce:     sync.Once{},
		gauges:       map[string]*Gauge{},
		histograms:   map[string]*Histogram{},
		mutex:        sync.Mutex{},
		shutdownOnce: sync.Once{},

		// Settings.
		httpEndpoint: config.HTTPEndpoint,
		httpHandler:  config.HTTPHandler,
		prefixes:     config.Prefixes,
	}

	return newService, nil
}

type Service struct {
	// Internals.
	closer       chan struct{}
	counters     map[string]*Counter
	bootOnce     sync.Once
	gauges       map[string]*Gauge
	histograms   map[string]*Histogram
	mutex        sync.Mutex
	shutdownOnce sync.Once

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

func (s *Service) Boot() {
	s.bootOnce.Do(func() {
		// Service specific boot logic goes here.
	})
}

func (s *Service) Counter(config spec.CounterConfig) (spec.Counter, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if c, ok := s.counters[config.Name()]; ok {
		return c, nil
	}

	newCounter, err := NewCounter(config)
	if err != nil {
		return nil, maskAny(err)
	}

	err = prometheus.Register(newCounter.ClientCounter)
	if err != nil {
		return nil, maskAny(err)
	}
	s.counters[config.Name()] = newCounter

	return newCounter, nil
}

func (s *Service) CounterConfig() spec.CounterConfig {
	return DefaultCounterConfig()
}

func (s *Service) Gauge(config spec.GaugeConfig) (spec.Gauge, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if g, ok := s.gauges[config.Name()]; ok {
		return g, nil
	}

	newGauge, err := NewGauge(config)
	if err != nil {
		return nil, maskAny(err)
	}

	err = prometheus.Register(newGauge.ClientGauge)
	if err != nil {
		return nil, maskAny(err)
	}
	s.gauges[config.Name()] = newGauge

	return newGauge, nil
}

func (s *Service) GaugeConfig() spec.GaugeConfig {
	return DefaultGaugeConfig()
}

func (s *Service) Histogram(config spec.HistogramConfig) (spec.Histogram, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if h, ok := s.histograms[config.Name()]; ok {
		return h, nil
	}

	newHistogram, err := NewHistogram(config)
	if err != nil {
		return nil, maskAny(err)
	}

	err = prometheus.Register(newHistogram.ClientHistogram)
	if err != nil {
		return nil, maskAny(err)
	}
	s.histograms[config.Name()] = newHistogram

	return newHistogram, nil
}

func (s *Service) HistogramConfig() spec.HistogramConfig {
	return DefaultHistogramConfig()
}

func (s *Service) HTTPEndpoint() string {
	return s.httpEndpoint
}

func (s *Service) HTTPHandler() http.Handler {
	return s.httpHandler
}

func (s *Service) Prefixes() []string {
	return s.prefixes
}

func (s *Service) NewKey(str ...string) string {
	return strings.Join(append(s.prefixes, str...), "_")
}

func (s *Service) Shutdown() {
	s.shutdownOnce.Do(func() {
		close(s.closer)
	})
}

func (s *Service) WrapFunc(key string, action func() error) func() error {
	wrappedFunc := func() error {
		histogramConfig := DefaultHistogramConfig()
		histogramConfig.SetName(s.NewKey(key, "milliseconds"))
		h, err := s.Histogram(histogramConfig)
		if err != nil {
			return maskAny(err)
		}
		counterConfig := DefaultCounterConfig()
		counterConfig.SetName(s.NewKey(key, "error", "total"))
		c, err := s.Counter(counterConfig)
		if err != nil {
			return maskAny(err)
		}

		defer func(t time.Time) {
			h.Observe(float64(time.Since(t) / time.Millisecond))
		}(time.Now())

		err = action()
		if err != nil {
			c.Increment(1)
			return maskAny(err)
		}

		return nil
	}

	return wrappedFunc
}
