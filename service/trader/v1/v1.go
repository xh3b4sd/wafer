package v1

import (
	"fmt"
	"reflect"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/xh3b4sd/wafer/service/buyer"
	"github.com/xh3b4sd/wafer/service/client"
	"github.com/xh3b4sd/wafer/service/informer"
	"github.com/xh3b4sd/wafer/service/seller"
	"github.com/xh3b4sd/wafer/service/trader"
	"github.com/xh3b4sd/wafer/service/trader/runtime"
	"github.com/xh3b4sd/wafer/service/trader/runtime/config"
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

	// Settings.
	Runtime config.Config
}

// DefaultConfig returns the default configuration used to create a new trader
// by best effort.
func DefaultConfig() Config {
	runtimeConfig := config.Config{}
	runtimeConfig.Trade.Budget = 500.0

	config := Config{
		// Dependencies.
		Buyer:    nil,
		Client:   nil,
		Informer: nil,
		Logger:   nil,
		Seller:   nil,

		// Settings.
		Runtime: runtimeConfig,
	}

	return config
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

	// Settings.
	err := config.Runtime.Validate()
	if err != nil {
		return nil, microerror.MaskAnyf(invalidConfigError, err.Error())
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
			Config: config.Runtime,
			State:  state.State{},
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
	var buys []informer.Price

	informerPrices := t.informer.Prices()
	t.runtime.State.Trade.Cycles = make([]int64, len(informerPrices))
	t.runtime.State.Trade.Revenues = make([]float64, len(informerPrices))

	for i, c := range informerPrices {
		for p := range c {
			// Manage sell events.
			for _, b := range buys {
				isSell, err := t.seller.Sell(p, b)
				if err != nil {
					return microerror.MaskAny(err)
				}

				if !isSell {
					continue
				}
				v := calculateVolume(p.Sell, t.runtime.Config.Trade.Budget)
				err = t.client.Sell(p, v)
				if err != nil {
					return microerror.MaskAny(err)
				}
				t.logger.Log("event", "sell", "price", fmt.Sprintf("%.2f", p.Sell))

				t.runtime.State.Trade.Cycles[i]++
				t.runtime.State.Trade.Revenues[i] += (p.Sell * v) - (b.Buy * v)
				t.buyer.DecrTradeConcurrent()
				buys = removePrice(buys, b)
			}

			// Manage buy events.
			{
				isBuy, err := t.buyer.Buy(p)
				if err != nil {
					return microerror.MaskAny(err)
				}

				if !isBuy {
					continue
				}
				buys = append(buys, p)
				err = t.client.Buy(p, calculateVolume(p.Buy, t.runtime.Config.Trade.Budget))
				if err != nil {
					return microerror.MaskAny(err)
					continue
				}
				t.logger.Log("event", "buy", "price", fmt.Sprintf("%.2f", p.Buy))
			}
		}
	}

	return nil
}

func (t *Trader) Runtime() runtime.Runtime {
	return t.runtime
}

func removePrice(buys []informer.Price, price informer.Price) []informer.Price {
	var list []informer.Price

	for _, p := range buys {
		if reflect.DeepEqual(p, price) {
			continue
		}

		list = append(list, p)
	}

	return list
}
