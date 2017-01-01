package spec

// Consumer represents a service to abstract instrumentation libraries to fetch
// application metrics.
//
// TODO extend with useful functionality.
type Consumer interface {
	// Boot initializes and starts the whole service like booting a machine. The
	// call to Boot blocks until the service is completely initialized, so you
	// might want to call it in a separate goroutine.
	Boot()
	// Shutdown ends all processes of the service like shutting down a machine.
	// The call to Shutdown blocks until the service is completely shut down, so
	// you might want to call it in a separate goroutine.
	Shutdown()
}
