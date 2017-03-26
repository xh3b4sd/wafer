package endpoint

import (
	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/xh3b4sd/wafer/server/endpoint/render"
	"github.com/xh3b4sd/wafer/server/endpoint/version"
	"github.com/xh3b4sd/wafer/server/middleware"
	"github.com/xh3b4sd/wafer/service"
)

// Config represents the configuration used to create a endpoint.
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

	var renderEndpoint *render.Endpoint
	{
		renderConfig := render.DefaultConfig()
		renderConfig.Logger = config.Logger
		renderConfig.Middleware = config.Middleware
		renderConfig.Service = config.Service
		renderEndpoint, err = render.New(renderConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var versionEndpoint *version.Endpoint
	{
		versionConfig := version.DefaultConfig()
		versionConfig.Logger = config.Logger
		versionConfig.Middleware = config.Middleware
		versionConfig.Service = config.Service
		versionEndpoint, err = version.New(versionConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	newEndpoint := &Endpoint{
		Render:  renderEndpoint,
		Version: versionEndpoint,
	}

	return newEndpoint, nil
}

// Endpoint is the endpoint collection.
type Endpoint struct {
	Render  *render.Endpoint
	Version *version.Endpoint
}
