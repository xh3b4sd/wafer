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

func (t *Trader) Execute() {
	done := make(chan struct{}, 1)
	buyPrice := informer.Price{}
	isBuyEvent := true

	go func() {
		for {
			select {
			case p := <-t.buyer.Buy():
				err := t.client.Buy(p)
				if err != nil {
					t.logger.Log("error", fmt.Sprintf("%#v", err))
					continue
				}

				buyPrice = p
				isBuyEvent = false
			case p := <-t.seller.Sell():
				err := t.client.Sell(p)
				if err != nil {
					t.logger.Log("error", fmt.Sprintf("%#v", err))
					continue
				}

				t.runtime.State.Trade.Revenue.Total += buyPrice.Buy - p.Sell
				isBuyEvent = true
			case <-done:
				return
			}
		}
	}()

	for p := range t.informer.Prices() {
		if isBuyEvent {
			err := t.buyer.Consume(p)
			if err != nil {
				t.logger.Log("error", fmt.Sprintf("%#v", err))
				continue
			}
		} else {
			err := t.seller.Consume(buyPrice, p)
			if err != nil {
				t.logger.Log("error", fmt.Sprintf("%#v", err))
				continue
			}
		}
	}

	t.client.Close()
	done <- struct{}{}
}

func (t *Trader) Runtime() runtime.Runtime {
	return t.runtime
}
