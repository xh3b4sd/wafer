package v1

import (
	"fmt"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/xh3b4sd/wafer/service/buyer"
	"github.com/xh3b4sd/wafer/service/client"
	"github.com/xh3b4sd/wafer/service/informer"
	"github.com/xh3b4sd/wafer/service/seller"
	"github.com/xh3b4sd/wafer/service/trader"
	"github.com/xh3b4sd/wafer/service/trader/runtime"
	"github.com/xh3b4sd/wafer/service/trader/runtime/state"
)

// Config is the configuration used to create a new trader.
type Config struct {
	// Dependencies.
	Buyer    buyer.Buyer
	Client   client.Client
	Informer informer.Informer
	Logger   micrologger.Logger
	Seller   seller.Seller
}

// DefaultConfig returns the default configuration used to create a new trader
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Buyer:    nil,
		Client:   nil,
		Informer: nil,
		Logger:   nil,
		Seller:   nil,
	}
}

// New creates a new configured trader.
func New(config Config) (trader.Trader, error) {
	// Dependencies.
	if config.Buyer == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Buyer must not be empty")
	}
	if config.Client == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Client must not be empty")
	}
	if config.Informer == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Informer must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Seller == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Seller must not be empty")
	}

	newTrader := &Trader{
		// Dependencies.
		buyer:    config.Buyer,
		client:   config.Client,
		informer: config.Informer,
		logger:   config.Logger,
		seller:   config.Seller,

		// Internals.
		runtime: runtime.Runtime{
			State: state.State{},
		},
	}

	return newTrader, nil
}

// Trader implements trader.Trader.
type Trader struct {
	// Dependencies.
	buyer    buyer.Buyer
	client   client.Client
	informer informer.Informer
	logger   micrologger.Logger
	seller   seller.Seller

	// Internals.
	runtime runtime.Runtime
}

func (t *Trader) Execute() error {
	buyPrice := informer.Price{}
	isBuyEvent := true

	for _, c := range t.informer.Prices() {
		for p := range c {
			if isBuyEvent {
				isBuy, err := t.buyer.Buy(p)
				if err != nil {
					return microerror.MaskAny(err)
				}

				if !isBuy {
					continue
				}
				err = t.client.Buy(p)
				if err != nil {
					return microerror.MaskAny(err)
					continue
				}
				t.logger.Log("event", "buy", "price", fmt.Sprintf("%.2f", p.Buy))

				buyPrice = p
				isBuyEvent = false
			} else {
				isSell, err := t.seller.Sell(buyPrice, p)
				if err != nil {
					return microerror.MaskAny(err)
				}

				if !isSell {
					continue
				}

				err = t.client.Sell(p)
				if err != nil {
					return microerror.MaskAny(err)
					continue
				}
				t.logger.Log("event", "sell", "price", fmt.Sprintf("%.2f", p.Sell))

				t.runtime.State.Trade.Revenue.Total += p.Sell - buyPrice.Buy
				isBuyEvent = true
			}
		}
	}

	return nil
}

func (t *Trader) Runtime() runtime.Runtime {
	return t.runtime
}
