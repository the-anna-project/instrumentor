// Package publisher implements
// github.com/the-anna-project/instrumentor.Publisher and provides
// instrumentation primitives to emit application metrics.
package publisher

import (
	"net/http"
	"strings"
	"sync"

	"github.com/the-anna-project/instrumentor/spec"
)

// ServiceConfig represents the configuration used to create a new memory
// publisher service.
type ServiceConfig struct {
}

// DefaultServiceConfig provides a default configuration to create a new memory
// publisher service by best effort.
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{}
}

// NewService creates a new memory publisher service.
func NewService(config ServiceConfig) (*Service, error) {
	newService := &Service{
		// Internals.
		closer:       make(chan struct{}, 1),
		bootOnce:     sync.Once{},
		shutdownOnce: sync.Once{},
	}

	return newService, nil
}

type Service struct {
	// Internals.
	closer       chan struct{}
	bootOnce     sync.Once
	shutdownOnce sync.Once
}

func (s *Service) Boot() {
	s.bootOnce.Do(func() {
		// Service specific boot logic goes here.
	})
}

func (s *Service) Counter(config spec.CounterConfig) (spec.Counter, error) {
	newCounter, err := NewCounter(config)
	if err != nil {
		return nil, maskAny(err)
	}

	return newCounter, nil
}

func (s *Service) CounterConfig() spec.CounterConfig {
	return DefaultCounterConfig()
}

func (s *Service) Gauge(config spec.GaugeConfig) (spec.Gauge, error) {
	newGauge, err := NewGauge(config)
	if err != nil {
		return nil, maskAny(err)
	}

	return newGauge, nil
}

func (s *Service) GaugeConfig() spec.GaugeConfig {
	return DefaultGaugeConfig()
}

func (s *Service) Histogram(config spec.HistogramConfig) (spec.Histogram, error) {
	newHistogram, err := NewHistogram(config)
	if err != nil {
		return nil, maskAny(err)
	}

	return newHistogram, nil
}

func (s *Service) HistogramConfig() spec.HistogramConfig {
	return DefaultHistogramConfig()
}

func (s *Service) HTTPEndpoint() string {
	return ""
}

func (s *Service) HTTPHandler() http.Handler {
	return nil
}

func (s *Service) Prefixes() []string {
	return nil
}

func (s *Service) NewKey(str ...string) string {
	return strings.Join(str, "_")
}

func (s *Service) Shutdown() {
	s.shutdownOnce.Do(func() {
		close(s.closer)
	})
}

func (s *Service) WrapFunc(key string, action func() error) func() error {
	wrappedFunc := func() error {
		err := action()
		if err != nil {
			return maskAny(err)
		}

		return nil
	}

	return wrappedFunc
}
