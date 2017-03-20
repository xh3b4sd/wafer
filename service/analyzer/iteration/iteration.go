package iteration

import (
	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/xh3b4sd/wafer/service/analyzer"
	analyzerconfig "github.com/xh3b4sd/wafer/service/analyzer/iteration/config"
	"github.com/xh3b4sd/wafer/service/decider"
	performancedecider "github.com/xh3b4sd/wafer/service/decider/performance"
	"github.com/xh3b4sd/wafer/service/informer"
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
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

	newAnalyzer := &Analyzer{
		// Dependencies.
		informer: config.Informer,
		logger:   config.Logger,

		// Internals.
		config: analyzerconfig.Config{},
	}

	return newAnalyzer, nil
}

// Analyzer implements analyzer.Analyzer.
type Analyzer struct {
	// Dependencies.
	informer informer.Informer
	logger   micrologger.Logger

	// Internals.
	config analyzerconfig.Config
}

// TODO
func (a *Analyzer) Adjust(config decider.Config) (decider.Config, error) {
	modified := config

	return modified, nil
}

func (a *Analyzer) Config() decider.Config {
	return a.config.Decider.Config.Best
}

func (a *Analyzer) Iterate(config decider.Config) (decider.Config, error) {
	done := make(chan struct{}, 1)

	var err error
	var revenue float64
	var newDecider decider.Decider
	{
		newConfig := performancedecider.DefaultConfig()
		newConfig.DeciderConfig = config
		newConfig.Logger = a.logger
		newDecider, err = performancedecider.New(newConfig)
		if err != nil {
			return decider.Config{}, microerror.MaskAny(err)
		}
	}

	// listen to buy and sell events of decider
	go func() {
		var buyPrice float64

		for {
			select {
			case p := <-newDecider.Buy():
				buyPrice = p.Buy
				a.logger.Log("event", "buy", "price", p.Buy)
			case p := <-newDecider.Sell():
				revenue += p.Sell - buyPrice
				a.logger.Log("event", "sell", "price", p.Sell)
			case <-done:
				return
			}
		}
	}()

	newDecider.Watch(a.informer.Prices())
	done <- struct{}{}

	if a.config.Analyzer.Runtime.Revenue.Best < revenue {
		a.config.Analyzer.Runtime.Revenue.Best = revenue
		a.config.Decider.Config.Best = config
	}

	adjusted, err := a.Adjust(config)
	if err != nil {
		return decider.Config{}, microerror.MaskAny(err)
	}

	return adjusted, nil
}

func (a *Analyzer) Revenue() float64 {
	return a.config.Analyzer.Runtime.Revenue.Best
}
