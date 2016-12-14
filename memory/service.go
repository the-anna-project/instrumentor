// Package memory mocks github.com/the-anna-project/instrumentor.Service and
// does effectively nothing. It is only used for default configurations that
// require a satisfied instrumentation implementation. This should then be
// overwritten with a valid implementation if required.
package memory

import (
	"net/http"

	"github.com/the-anna-project/instrumentor"
)

// Config represents the configuration used to create a new instrumentor
// service.
type Config struct {
}

// DefaultConfig provides a default configuration to create a new instrumentor
// service by best effort.
func DefaultConfig() Config {
	return Config{}
}

// New creates a new instrumentor service.
func New(config Config) (instrumentor.Service, error) {
	newService := &service{}

	return newService, nil
}

type service struct {
}

func (s *service) ExecFunc(key string, action func() error) error {
	err := action()
	if err != nil {
		return maskAny(err)
	}

	return nil
}

func (s *service) Counter(key string) (instrumentor.Counter, error) {
	newConfig := DefaultCounterConfig()
	newCounter, err := NewCounter(newConfig)
	if err != nil {
		return nil, maskAny(err)
	}

	return newCounter, nil
}

func (s *service) Gauge(key string) (instrumentor.Gauge, error) {
	newConfig := DefaultGaugeConfig()
	newGauge, err := NewGauge(newConfig)
	if err != nil {
		return nil, maskAny(err)
	}

	return newGauge, nil
}

func (s *service) Histogram(key string) (instrumentor.Histogram, error) {
	newConfig := DefaultHistogramConfig()
	newHistogram, err := NewHistogram(newConfig)
	if err != nil {
		return nil, maskAny(err)
	}

	return newHistogram, nil
}

func (s *service) HTTPEndpoint() string {
	return ""
}

func (s *service) HTTPHandler() http.Handler {
	return nil
}

func (s *service) Prefixes() []string {
	return nil
}

func (s *service) NewKey(str ...string) string {
	return ""
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
