package v1

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/xh3b4sd/wafer/service/analyzer"
	"github.com/xh3b4sd/wafer/service/analyzer/runtime"
	runtimeconfig "github.com/xh3b4sd/wafer/service/analyzer/runtime/config"
	runtimestate "github.com/xh3b4sd/wafer/service/analyzer/runtime/state"
	statehistory "github.com/xh3b4sd/wafer/service/analyzer/runtime/state/config/history"
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
	Informer informer.Informer
	Logger   micrologger.Logger
}

// DefaultConfig returns the default configuration used to create a new analyzer
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Informer: nil,
		Logger:   nil,
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

	runtimeConfig := &runtimeconfig.Config{}

	var newPermutation permutation.Permutation
	var err error
	{
		config := v1permutation.DefaultConfig()
		config.Logger = config.Logger
		config.Object = runtimeConfig
		newPermutation, err = v1permutation.New(config)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	newAnalyzer := &Analyzer{
		// Dependencies.
		informer: config.Informer,
		logger:   config.Logger,

		// Internals.
		analyzeOnce: sync.Once{},
		mutex:       sync.Mutex{},
		permutation: newPermutation,
		runtime: runtime.Runtime{
			Config: runtimeConfig,
			State:  runtimestate.State{},
		},
	}

	return newAnalyzer, nil
}

// Analyzer implements analyzer.Analyzer.
type Analyzer struct {
	// Dependencies.
	informer informer.Informer
	logger   micrologger.Logger

	// Internals.
	analyzeOnce sync.Once
	mutex       sync.Mutex
	permutation permutation.Permutation
	runtime     runtime.Runtime
}

func (a *Analyzer) Execute() {
	a.analyzeOnce.Do(func() {
		err := a.execute()
		if err != nil {
			a.logger.Log("error", fmt.Sprintf("%#v", err))
		}
	})
}

func (a *Analyzer) Runtime() runtime.Runtime {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.runtime
}

func (a *Analyzer) execute() error {
	indizes := v1permutation.IndizesFromConfigs(a.runtime.Config.GetPermConfigs())
	zeroIndizes := v1permutation.IndizesFromConfigs(a.runtime.Config.GetPermConfigs())
	max := v1permutation.MaxFromConfigs(a.runtime.Config.GetPermConfigs())
	stepDuration := &Duration{}

	var stepCurrent float64
	a.runtime.State.Informer.Prices = a.informer.Runtime().State.Prices
	a.runtime.State.Permutation.Max = max
	a.runtime.State.Permutation.Start = time.Now()
	a.runtime.State.Permutation.Step.Total = v1permutation.TotalFromMax(max)

	for {
		var stepStart time.Time
		{
			stepStart = time.Now()
			stepCurrent++
			a.mutex.Lock()
			a.runtime.State.Permutation.Step.Current = stepCurrent
			a.runtime.State.Permutation.Indizes = indizes
			a.mutex.Unlock()
		}

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

			// We have to set some parts of the configuration to the runtime config of
			// the analyzer to be able to track the configuration history properly for
			// the permutation process.
			runtimeConfig.Buyer.Trade.Concurrent = config.Runtime.Trade.Concurrent
			config.Runtime = runtimeConfig.Buyer

			newBuyer, err = v1buyer.New(config)
			if err != nil {
				return microerror.MaskAny(err)
			}
		}

		var newClient client.Client
		{
			config := analyzerclient.DefaultConfig()
			config.Logger = a.logger
			newClient, err = analyzerclient.New(config)
			if err != nil {
				return microerror.MaskAny(err)
			}
		}

		var newSeller seller.Seller
		{
			config := v1seller.DefaultConfig()
			config.Logger = a.logger

			// We have to set some parts of the configuration to the runtime config of
			// the analyzer to be able to track the configuration history properly for
			// the permutation process.
			runtimeConfig.Seller.Trade.Fee.Min = config.Runtime.Trade.Fee.Min
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

		a.mutex.Lock()
		cycles := newTrader.Runtime().State.Trade.Cycles
		revenues := newTrader.Runtime().State.Trade.Revenues
		if (len(a.runtime.State.Config.History) == 0 && sum(revenues) > 0) || (len(a.runtime.State.Config.History) > 0 && sum(a.runtime.State.Config.History[0].Revenues) < sum(revenues)) {
			history := statehistory.History{
				Config:   *runtimeConfig,
				Cycles:   cycles,
				Indizes:  append([]int{}, indizes...), // copy
				Revenues: revenues,
			}
			a.runtime.State.Config.History = append([]statehistory.History{history}, a.runtime.State.Config.History...) // prepend
		}
		a.mutex.Unlock()

		indizes, err = v1permutation.ShiftIndizes(indizes, max)
		if err != nil {
			return microerror.MaskAny(err)
		}

		{
			stepEnd := time.Now()
			stepDuration.Add(stepEnd.Sub(stepStart))
			a.mutex.Lock()
			a.runtime.State.Permutation.Progress = fmt.Sprintf("%.3f", a.runtime.State.Permutation.Step.Current*100/a.runtime.State.Permutation.Step.Total)
			a.runtime.State.Permutation.Step.Duration = fmt.Sprintf("%.3fs", stepDuration.Average().Seconds())
			a.runtime.State.Permutation.End = a.eta(stepDuration.Average())
			a.mutex.Unlock()
		}

		if reflect.DeepEqual(indizes, zeroIndizes) {
			break
		}
	}

	return nil
}

func (a *Analyzer) eta(stepDuration time.Duration) time.Time {
	dur := time.Duration(a.runtime.State.Permutation.Step.Total-a.runtime.State.Permutation.Step.Current) * stepDuration
	eta := time.Now().Add(dur)

	return eta
}

func sum(list []float64) float64 {
	var s float64

	for _, f := range list {
		s += f
	}

	return s
}
