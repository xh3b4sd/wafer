package v1

import (
	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/xh3b4sd/wafer/service/informer"
	"github.com/xh3b4sd/wafer/service/seller"
	"github.com/xh3b4sd/wafer/service/seller/runtime"
	"github.com/xh3b4sd/wafer/service/seller/runtime/config"
	"github.com/xh3b4sd/wafer/service/seller/runtime/state"
)

// Config is the configuration used to create a new seller.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Runtime config.Config
}

// DefaultConfig returns the default configuration used to create a new seller
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Runtime: config.Config{},
	}
}

// New creates a new configured seller.
func New(config Config) (seller.Seller, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Logger must not be empty")
	}

	newSeller := &Seller{
		// Dependencies.
		logger: config.Logger,

		// Internals.
		runtime: runtime.Runtime{
			Config: config.Runtime,
			State:  state.State{},
		},
	}

	return newSeller, nil
}

// Seller implements seller.Seller.
type Seller struct {
	// Dependencies.
	logger micrologger.Logger

	// Internals.
	runtime runtime.Runtime
}

func (s *Seller) Runtime() runtime.Runtime {
	return s.runtime
}

func (s *Seller) Sell(buyPrice, currentPrice informer.Price) (bool, error) {
	revenue := calculateRevenue(buyPrice.Buy, currentPrice.Sell, s.runtime.Config.Trade.Fee.Min)
	if revenue < s.runtime.Config.Trade.Revenue.Min {
		return false, nil
	}

	duration := currentPrice.Time.Sub(buyPrice.Time)
	if duration < s.runtime.Config.Trade.Duration.Min {
		return false, nil
	}

	return true, nil
}
