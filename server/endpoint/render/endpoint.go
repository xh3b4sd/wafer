package render

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	kitendpoint "github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"

	"github.com/xh3b4sd/wafer/server/middleware"
	"github.com/xh3b4sd/wafer/service"
	"github.com/xh3b4sd/wafer/service/buyer"
	v1buyer "github.com/xh3b4sd/wafer/service/buyer/v1"
	"github.com/xh3b4sd/wafer/service/client"
	analyzerclient "github.com/xh3b4sd/wafer/service/client/analyzer"
	"github.com/xh3b4sd/wafer/service/informer"
	"github.com/xh3b4sd/wafer/service/informer/csv"
	"github.com/xh3b4sd/wafer/service/seller"
	v1seller "github.com/xh3b4sd/wafer/service/seller/v1"
	"github.com/xh3b4sd/wafer/service/trader"
	v1trader "github.com/xh3b4sd/wafer/service/trader/v1"
)

const (
	// Method is the HTTP method this endpoint is registered for.
	Method = "GET"
	// Name identifies the endpoint. It is aligned to the package path.
	Name = "render"
	// Path is the HTTP request path this endpoint is registered for.
	Path = "/v1/render/"
)

// Config represents the configuration used to create a version endpoint.
type Config struct {
	// Dependencies.
	Logger     micrologger.Logger
	Middleware *middleware.Middleware
	Service    *service.Service
}

// DefaultConfig provides a default configuration to create a new version
// endpoint by best effort.
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
		w.Header().Set(http.CanonicalHeaderKey("Content-Type"), "image/png")

		w.WriteHeader(http.StatusOK)

		_, err := w.Write(response.(*bytes.Buffer).Bytes())
		if err != nil {
			return microerror.MaskAny(err)
		}

		return nil
	}
}

func (e *Endpoint) Endpoint() kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var err error

		buyChan := make(chan informer.Price, 100)
		sellChan := make(chan informer.Price, 100)

		var newInformer informer.Informer
		{
			config := csv.DefaultConfig()
			config.File.Header.Buy = 9
			config.File.Header.Ignore = true
			config.File.Header.Sell = 10
			config.File.Header.Time = 12
			config.File.Path = "/Users/xh3b4sd/Downloads/test.csv"
			newInformer, err = csv.New(config)
			if err != nil {
				return nil, microerror.MaskAny(err)
			}
		}

		var newClient client.Client
		{
			config := analyzerclient.DefaultConfig()
			config.BuyChan = buyChan
			config.Logger = e.logger
			config.SellChan = sellChan
			newClient, err = analyzerclient.New(config)
			if err != nil {
				return nil, microerror.MaskAny(err)
			}
			defer newClient.Close()
		}

		var newBuyer buyer.Buyer
		{
			config := v1buyer.DefaultConfig()
			config.Logger = e.logger
			newBuyer, err = v1buyer.New(config)
			if err != nil {
				return nil, microerror.MaskAny(err)
			}
		}

		var newSeller seller.Seller
		{
			config := v1seller.DefaultConfig()
			config.Logger = e.logger
			config.Runtime.Trade.Duration.Min = 4 * time.Minute
			config.Runtime.Trade.Fee.Min = 1.5
			config.Runtime.Trade.Revenue.Min = 4.5
			newSeller, err = v1seller.New(config)
			if err != nil {
				return nil, microerror.MaskAny(err)
			}
		}

		var newTrader trader.Trader
		{
			config := v1trader.DefaultConfig()
			config.Buyer = newBuyer
			config.Client = newClient
			config.Informer = newInformer
			config.Logger = e.logger
			config.Seller = newSeller
			newTrader, err = v1trader.New(config)
			if err != nil {
				return nil, microerror.MaskAny(err)
			}
		}

		//		var newExchange exchange.Exchange
		//		{
		//			config := analyzerexchange.DefaultConfig()
		//			config.BuyChan = buyChan
		//			config.Informer = newInformer
		//			config.Logger = e.logger
		//			config.SellChan = sellChan
		//			newExchange, err = analyzerexchange.New(config)
		//			if err != nil {
		//				return nil, microerror.MaskAny(err)
		//			}
		//			defer newExchange.Close()
		//		}

		//		go func() {
		//			e.logger.Log("debug", "exchange started")
		//			newExchange.Execute()
		//			e.logger.Log("debug", "exchange stopped")
		//		}()

		e.logger.Log("debug", "trader started")
		err = newTrader.Execute()
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
		e.logger.Log("debug", "trader stopped")

		//		e.logger.Log("debug", "rendering started")
		//		buf, err := newExchange.Render()
		//		if err != nil {
		//			return nil, microerror.MaskAny(err)
		//		}
		//		e.logger.Log("debug", "rendering finished")

		e.logger.Log("debug", fmt.Sprintf("trader revenue: %#v", newTrader.Runtime().State.Trade.Revenues))

		return nil, nil
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
