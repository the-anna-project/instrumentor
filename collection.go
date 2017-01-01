package instrumentor

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	memoryconsumer "github.com/the-anna-project/instrumentor/memory/consumer"
	memorypublisher "github.com/the-anna-project/instrumentor/memory/publisher"
	prometheusconsumer "github.com/the-anna-project/instrumentor/prometheus/consumer"
	prometheuspublisher "github.com/the-anna-project/instrumentor/prometheus/publisher"
	"github.com/the-anna-project/instrumentor/spec"
)

const (
	// KindMemory is the kind to be used to create a memory instrumentor services.
	KindMemory = "memory"
	// KindPrometheus is the kind to be used to create a collection of prometheus
	// instrumentor services.
	KindPrometheus = "prometheus"
)

// Config represents the configuration used to create a new collection.
type Config struct {
	// Settings.
	HTTPEndpoint string
	HTTPHandler  http.Handler
	Kind         string
	Prefixes     []string
}

// DefaultConfig provides a default configuration to create a new collection by
// best effort.
func DefaultConfig() Config {
	return Config{
		// Settings.
		HTTPEndpoint: "/metrics",
		HTTPHandler:  prometheus.Handler(),
		Kind:         KindMemory,
		Prefixes:     []string{},
	}
}

// New creates a new configured storage Collection.
func New(config Config) (*Collection, error) {
	// Settings.
	if config.Kind == "" {
		return nil, maskAnyf(invalidConfigError, "kind must not be empty")
	}
	if config.Kind != KindMemory && config.Kind != KindPrometheus {
		return nil, maskAnyf(invalidConfigError, "kind must be one of: %s, %s", KindMemory, KindPrometheus)
	}

	var err error

	var consumerService spec.Consumer
	{
		switch config.Kind {
		case KindMemory:
			consumerConfig := memoryconsumer.DefaultServiceConfig()
			consumerService, err = memoryconsumer.NewService(consumerConfig)
			if err != nil {
				return nil, maskAny(err)
			}
		case KindPrometheus:
			consumerConfig := prometheusconsumer.DefaultServiceConfig()
			consumerService, err = prometheusconsumer.NewService(consumerConfig)
			if err != nil {
				return nil, maskAny(err)
			}
		}
	}

	var publisherService spec.Publisher
	{
		switch config.Kind {
		case KindMemory:
			publisherConfig := memorypublisher.DefaultServiceConfig()
			publisherService, err = memorypublisher.NewService(publisherConfig)
			if err != nil {
				return nil, maskAny(err)
			}
		case KindPrometheus:
			publisherConfig := prometheuspublisher.DefaultServiceConfig()
			publisherConfig.HTTPEndpoint = config.HTTPEndpoint
			publisherConfig.HTTPHandler = config.HTTPHandler
			publisherConfig.Prefixes = config.Prefixes
			publisherService, err = prometheuspublisher.NewService(publisherConfig)
			if err != nil {
				return nil, maskAny(err)
			}
		}
	}

	newCollection := &Collection{
		// Internals.
		bootOnce:     sync.Once{},
		shutdownOnce: sync.Once{},

		// Public.
		Consumer:  consumerService,
		Publisher: publisherService,
	}

	return newCollection, nil
}

// Collection is the object bundling all services.
type Collection struct {
	// Internals.
	bootOnce     sync.Once
	shutdownOnce sync.Once

	// Public.
	Consumer  spec.Consumer
	Publisher spec.Publisher
}

func (c *Collection) Boot() {
	c.bootOnce.Do(func() {
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			c.Consumer.Boot()
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			c.Publisher.Boot()
			wg.Done()
		}()

		wg.Wait()
	})
}

func (c *Collection) Shutdown() {
	c.shutdownOnce.Do(func() {
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			c.Consumer.Shutdown()
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			c.Publisher.Shutdown()
			wg.Done()
		}()

		wg.Wait()
	})
}
