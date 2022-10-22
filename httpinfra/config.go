package httpinfra

import "time"

type config struct {
	Host              string
	Port              int
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	TracerEnabled     bool
	TracerServiceName string
	Version           string
	Revision          string
}

// Option specifies server configuration options.
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

// WithHost sets the host for the server.
func WithHost(host string) Option {
	return optionFunc(func(c *config) {
		c.Host = host
	})
}

// WithPort sets the port for the server.
func WithPort(port int) Option {
	return optionFunc(func(c *config) {
		c.Port = port
	})
}

// WithReadTimeout sets the read timeout for the server.
func WithReadTimeout(timeout time.Duration) Option {
	return optionFunc(func(c *config) {
		c.ReadTimeout = timeout
	})
}

// WithWriteTimeout sets to write timeout for the server.
func WithWriteTimeout(timeout time.Duration) Option {
	return optionFunc(func(c *config) {
		c.WriteTimeout = timeout
	})
}

// WithIdleTimeout sets the idle timeout for the server.
func WithIdleTimeout(timeout time.Duration) Option {
	return optionFunc(func(c *config) {
		c.IdleTimeout = timeout
	})
}

// WithTracer enables the tracer for the server.
func WithTracer() Option {
	return optionFunc(func(c *config) {
		c.TracerEnabled = true
	})
}

// WithOtelServiceName sets the otel service name for the server.
func WithOtelServiceName(name string) Option {
	return optionFunc(func(c *config) {
		c.TracerServiceName = name
	})
}

// WithVersion sets the version for the server.
func WithVersion(version string) Option {
	return optionFunc(func(c *config) {
		c.Version = version
	})
}

// WithRevision sets the revision for the server.
func WithRevision(revision string) Option {
	return optionFunc(func(c *config) {
		c.Revision = revision
	})
}

// ApplyOptions applies the given options to the config.
func (c *config) ApplyOptions(opts ...Option) {
	for _, opt := range opts {
		opt.apply(c)
	}
}

// NewDefaultConfig returns a new default config.
func NewDefaultConfig() *config {
	return &config{
		Host:              "",
		Port:              8081,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       20 * time.Second,
		TracerEnabled:     false,
		TracerServiceName: "unknown_service",
		Version:           "unknown",
		Revision:          "unknown",
	}
}
