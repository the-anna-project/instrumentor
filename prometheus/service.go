// Package prometheus implements spec.InstrumentorService and provides instrumentation
// primitives to manage application metrics.
package prometheus

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	objectspec "github.com/the-anna-project/spec/object"
	servicespec "github.com/the-anna-project/spec/service"
)

// New creates a new pronetheus instrumentor service.
func New() servicespec.InstrumentorService {
	return &service{}
}

type service struct {
	// Dependencies.

	serviceCollection servicespec.ServiceCollection

	// Settings.

	counters   map[string]objectspec.InstrumentorCounter
	gauges     map[string]objectspec.InstrumentorGauge
	histograms map[string]objectspec.InstrumentorHistogram
	// httpEndpoint represents the HTTP endpoint used to register the
	// httpHandler. In the context of Prometheus this is usually /metrics.
	httpEndpoint string
	// httpHandler represents the HTTP handler used to register the Prometheus
	// registry in the HTTP server.
	httpHandler http.Handler
	metadata    map[string]string
	mutex       sync.Mutex
	// prefixes represents the Instrumentor's ordered prefixes. It is recommended
	// to use the following prefixes.
	//
	//     []string{"anna", "<prefix>"}
	//
	prefixes []string
}

func (s *service) Boot() {
	id, err := s.Service().ID().New()
	if err != nil {
		panic(err)
	}
	s.metadata = map[string]string{
		"id":   id,
		"kind": "prometheus",
		"name": "instrumentor",
		"type": "service",
	}

	s.counters = map[string]objectspec.InstrumentorCounter{}
	s.gauges = map[string]objectspec.InstrumentorGauge{}
	s.histograms = map[string]objectspec.InstrumentorHistogram{}
	s.httpEndpoint = "/metrics"
	s.httpHandler = prometheus.Handler()
	s.mutex = sync.Mutex{}
	s.prefixes = []string{"anna"}
}

func (s *service) ExecFunc(key string, action func() error) error {
	h, err := s.GetHistogram(s.NewKey(key, "durations", "histogram", "milliseconds"))
	if err != nil {
		return maskAny(err)
	}
	c, err := s.GetCounter(s.NewKey(key, "errors", "counter", "total"))
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

func (s *service) GetCounter(key string) (objectspec.InstrumentorCounter, error) {
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

func (s *service) GetGauge(key string) (objectspec.InstrumentorGauge, error) {
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

func (s *service) GetHistogram(key string) (objectspec.InstrumentorHistogram, error) {
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

func (s *service) GetHTTPEndpoint() string {
	return s.httpEndpoint
}

func (s *service) GetHTTPHandler() http.Handler {
	return s.httpHandler
}

func (s *service) GetPrefixes() []string {
	return s.prefixes
}

func (s *service) Metadata() map[string]string {
	return s.metadata
}

func (s *service) NewKey(str ...string) string {
	return strings.Join(append(s.prefixes, str...), "_")
}

func (s *service) Service() servicespec.ServiceCollection {
	return s.serviceCollection
}

func (s *service) SetServiceCollection(sc servicespec.ServiceCollection) {
	s.serviceCollection = sc
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
