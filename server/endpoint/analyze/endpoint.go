package analyze

import (
	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/xh3b4sd/wafer/server/endpoint/analyze/create"
	"github.com/xh3b4sd/wafer/server/endpoint/analyze/search"
	"github.com/xh3b4sd/wafer/server/middleware"
	"github.com/xh3b4sd/wafer/service"
)

// Config represents the configuration used to create an endpoint.
type Config struct {
	// Dependencies.
	Logger     micrologger.Logger
	Middleware *middleware.Middleware
	Service    *service.Service
}

// DefaultConfig provides a default configuration to create a new endpoint by
// best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger:     nil,
		Middleware: nil,
		Service:    nil,
	}
}

// New creates a new configured endpoint.
func New(config Config) (*Endpoint, error) {
	var err error

	var createEndpoint *create.Endpoint
	{
		createConfig := create.DefaultConfig()
		createConfig.Logger = config.Logger
		createConfig.Middleware = config.Middleware
		createConfig.Service = config.Service
		createEndpoint, err = create.New(createConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var searchEndpoint *search.Endpoint
	{
		searchConfig := search.DefaultConfig()
		searchConfig.Logger = config.Logger
		searchConfig.Middleware = config.Middleware
		searchConfig.Service = config.Service
		searchEndpoint, err = search.New(searchConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	newEndpoint := &Endpoint{
		Create: createEndpoint,
		Search: searchEndpoint,
	}

	return newEndpoint, nil
}

// Endpoint is the endpoint collection.
type Endpoint struct {
	Create *create.Endpoint
	Search *search.Endpoint
}
