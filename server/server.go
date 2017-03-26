// Package server provides a server implementation to connect network transport
// protocols and service business logic by defining server endpoints.
package server

import (
	"net/http"
	"sync"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	microserver "github.com/giantswarm/microkit/server"
	kithttp "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"

	"github.com/xh3b4sd/wafer/server/endpoint"
	"github.com/xh3b4sd/wafer/server/middleware"
	"github.com/xh3b4sd/wafer/service"
)

// Config represents the configuration used to create a new server object.
type Config struct {
	// Dependencies.
	Service *service.Service

	// Settings.
	MicroServerConfig microserver.Config
}

// DefaultConfig provides a default configuration to create a new server object
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Service: nil,

		// Settings.
		MicroServerConfig: microserver.DefaultConfig(),
	}
}

// New creates a new configured server object.
func New(config Config) (microserver.Server, error) {
	var err error

	var middlewareCollection *middleware.Middleware
	{
		middlewareConfig := middleware.DefaultConfig()
		middlewareConfig.Logger = config.MicroServerConfig.Logger
		middlewareConfig.Service = config.Service
		middlewareCollection, err = middleware.New(middlewareConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	var endpointCollection *endpoint.Endpoint
	{
		endpointConfig := endpoint.DefaultConfig()
		endpointConfig.Logger = config.MicroServerConfig.Logger
		endpointConfig.Middleware = middlewareCollection
		endpointConfig.Service = config.Service
		endpointCollection, err = endpoint.New(endpointConfig)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	newServer := &server{
		// Dependencies.
		logger: config.MicroServerConfig.Logger,

		// Internals.
		bootOnce:     sync.Once{},
		config:       config.MicroServerConfig,
		shutdownOnce: sync.Once{},
	}

	// Apply internals to the micro server config.
	newServer.config.Endpoints = []microserver.Endpoint{
		endpointCollection.Render,
		endpointCollection.Version,
	}
	newServer.config.ErrorEncoder = newServer.newErrorEncoder()

	return newServer, nil
}

type server struct {
	// Dependencies.
	logger micrologger.Logger

	// Internals.
	bootOnce     sync.Once
	config       microserver.Config
	shutdownOnce sync.Once
}

func (s *server) Boot() {
	s.bootOnce.Do(func() {
		// Here goes your custom boot logic for your server/endpoint/middleware, if
		// any.
	})
}

func (s *server) Config() microserver.Config {
	return s.config
}

func (s *server) Shutdown() {
	s.shutdownOnce.Do(func() {
		// Here goes your custom shutdown logic for your server/endpoint/middleware,
		// if any.
	})
}

func (s *server) newErrorEncoder() kithttp.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		rErr := err.(microserver.ResponseError)
		//uErr := rErr.Underlying()

		if rErr.IsEndpoint() {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}