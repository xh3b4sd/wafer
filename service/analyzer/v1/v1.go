package v1

import (
	"reflect"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/xh3b4sd/wafer/service/analyzer"
	"github.com/xh3b4sd/wafer/service/analyzer/runtime"
	runtimeconfig "github.com/xh3b4sd/wafer/service/analyzer/runtime/config"
	runtimestate "github.com/xh3b4sd/wafer/service/analyzer/runtime/state"
	statehistory "github.com/xh3b4sd/wafer/service/analyzer/runtime/state/history"
	"github.com/xh3b4sd/wafer/service/buyer"
	v1buyer "github.com/xh3b4sd/wafer/service/buyer/v1"
	"github.com/xh3b4sd/wafer/service/client"
	analyzerclient "github.com/xh3b4sd/wafer/service/client/analyzer"
	"github.com/xh3b4sd/wafer/service/informer"
	"github.com/xh3b4sd/wafer/service/permutation"
	v1permutation "github.com/xh3b4sd/wafer/service/permutation/v1"
	"github.com/xh3b4sd/wafer/service/seller"
	v1seller "github.com/xh3b4sd/wafer/service/seller/v1"
	"github.com/xh3b4sd/wafer/service/trader"
	v1trader "github.com/xh3b4sd/wafer/service/trader/v1"
)

// Config is the configuration used to create a new analyzer.
type Config struct {
	// Dependencies.
	Informer    informer.Informer
	Logger      micrologger.Logger
	Permutation permutation.Permutation
}

// DefaultConfig returns the default configuration used to create a new analyzer
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Informer:    nil,
		Logger:      nil,
		Permutation: nil,
	}
}

// New creates a new configured analyzer.
func New(config Config) (analyzer.Analyzer, error) {
	// Dependencies.
	if config.Informer == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Informer must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Logger must not be empty")
	}
	if config.Permutation == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Permutation must not be empty")
	}

	newAnalyzer := &Analyzer{
		// Dependencies.
		informer:    config.Informer,
		logger:      config.Logger,
		permutation: config.Permutation,

		// Internals.
		runtime: runtime.Runtime{
			Config: &runtimeconfig.Config{},
			State:  runtimestate.State{},
		},

		// Settings.
		buys:  []informer.Price{},
		sells: []informer.Price{},
	}

	// Initilized empty history record with zero values. The algorithm below
	// requires this initialization.
	newAnalyzer.runtime.State.Config.History = append(newAnalyzer.runtime.State.Config.History, statehistory.History{})

	return newAnalyzer, nil
}

// Analyzer implements analyzer.Analyzer.
type Analyzer struct {
	// Dependencies.
	informer    informer.Informer
	logger      micrologger.Logger
	permutation permutation.Permutation

	// Internals.
	runtime runtime.Runtime

	// Settings.
	buys  []informer.Price
	sells []informer.Price
}

func (a *Analyzer) Execute() error {
	indizes := v1permutation.IndizesFromConfigs(a.runtime.Config.GetPermConfigs())
	zeroIndizes := v1permutation.IndizesFromConfigs(a.runtime.Config.GetPermConfigs())
	max := v1permutation.MaxFromConfigs(a.runtime.Config.GetPermConfigs())

	for {
		permConfig, err := a.permutation.ValueFor(indizes)
		if err != nil {
			return microerror.MaskAny(err)
		}
		runtimeConfig, ok := permConfig.(*runtimeconfig.Config)
		if !ok {
			return microerror.MaskAnyf(invalidExecutionError, "invalid type for runtime config")
		}

		var newBuyer buyer.Buyer
		{
			config := v1buyer.DefaultConfig()
			config.Logger = a.logger
			config.Runtime = runtimeConfig.Buyer
			newBuyer, err = v1buyer.New(config)
			if err != nil {
				return microerror.MaskAny(err)
			}
		}

		var newClient client.Client
		{
			config := analyzerclient.DefaultConfig()
			config.Discard = true
			config.Logger = a.logger
			newClient, err = analyzerclient.New(config)
			if err != nil {
				return microerror.MaskAny(err)
			}
			defer newClient.Close()
		}

		var newSeller seller.Seller
		{
			config := v1seller.DefaultConfig()
			config.Logger = a.logger
			config.Runtime = runtimeConfig.Seller
			newSeller, err = v1seller.New(config)
			if err != nil {
				return microerror.MaskAny(err)
			}
		}

		var newTrader trader.Trader
		{
			config := v1trader.DefaultConfig()
			config.Buyer = newBuyer
			config.Client = newClient
			config.Informer = a.informer
			config.Logger = a.logger
			config.Runtime = runtimeConfig.Trader
			config.Seller = newSeller
			newTrader, err = v1trader.New(config)
			if err != nil {
				return microerror.MaskAny(err)
			}
		}

		err = newTrader.Execute()
		if err != nil {
			return microerror.MaskAny(err)
		}

		revenue := newTrader.Runtime().State.Trade.Revenue
		if a.runtime.State.Config.History[0].Revenue < revenue {
			history := statehistory.History{
				Config:  *runtimeConfig,
				Indizes: append([]int{}, indizes...), // copy
				Revenue: revenue,
			}
			a.runtime.State.Config.History = append([]statehistory.History{history}, a.runtime.State.Config.History...) // prepend
		}

		indizes, err = v1permutation.ShiftIndizes(indizes, max)
		if err != nil {
			return microerror.MaskAny(err)
		}

		if reflect.DeepEqual(indizes, zeroIndizes) {
			break
		}
	}

	// TODO track duration of iterations in runtime state

	// TODO track calculated execution progress in runtime state

	return nil
}

func (a *Analyzer) Runtime() runtime.Runtime {
	return a.runtime
}
