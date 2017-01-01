// Package consumer implements github.com/the-anna-project/instrumentor.Consumer
// and provides instrumentation primitives to manage application metrics.
package consumer

import (
	"sync"
)

// ServiceConfig represents the configuration used to create a new prometheus
// consumer service.
type ServiceConfig struct {
}

// DefaultServiceConfig provides a default configuration to create a new
// prometheus consumer service by best effort.
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{}
}

// NewService creates a new prometheus consumer service.
func NewService(config ServiceConfig) (*Service, error) {
	newService := &Service{
		// Internals.
		bootOnce:     sync.Once{},
		closer:       make(chan struct{}, 1),
		shutdownOnce: sync.Once{},
	}

	return newService, nil
}

type Service struct {
	// Internals.
	bootOnce     sync.Once
	closer       chan struct{}
	shutdownOnce sync.Once
}

func (s *Service) Boot() {
	s.bootOnce.Do(func() {
		// Service specific boot logic goes here.
	})
}

func (s *Service) Shutdown() {
	s.shutdownOnce.Do(func() {
		close(s.closer)
	})
}
