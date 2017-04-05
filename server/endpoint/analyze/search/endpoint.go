package search

import (
	"encoding/json"
	"net/http"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"

	"github.com/xh3b4sd/wafer/server/middleware"
	"github.com/xh3b4sd/wafer/service"
)

const (
	// Method is the HTTP method this endpoint is registered for.
	Method = "GET"
	// Name identifies the endpoint. It is aligned to the package path.
	Name = "analyze/search"
	// Path is the HTTP request path this endpoint is registered for.
	Path = "/v1/analyze/{id}"
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

// New creates a new configured version endpoint.
func New(config Config) (*Endpoint, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Middleware == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Middleware must not be empty")
	}
	if config.Service == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Service must not be empty")
	}

	newEndpoint := &Endpoint{
		// Dependencies.
		logger:     config.Logger,
		middleware: config.Middleware,
		service:    config.Service,
	}

	return newEndpoint, nil
}

type Endpoint struct {
	// Dependencies.
	logger     micrologger.Logger
	middleware *middleware.Middleware
	service    *service.Service
}

func (e *Endpoint) Decoder() kithttp.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		return nil, nil
	}
}

func (e *Endpoint) Encoder() kithttp.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		w.Header().Set(http.CanonicalHeaderKey("Content-Type"), "application/json; charset=utf-8")

		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			return microerror.MaskAny(err)
		}

		return nil
	}
}

func (e *Endpoint) Endpoint() kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		response := Response{}

		response.Analyzer.State = e.service.Analyzer.Runtime().State

		return response, nil
	}
}

func (e *Endpoint) Method() string {
	return Method
}

func (e *Endpoint) Middlewares() []kitendpoint.Middleware {
	return []kitendpoint.Middleware{}
}

func (e *Endpoint) Name() string {
	return Name
}

func (e *Endpoint) Path() string {
	return Path
}
