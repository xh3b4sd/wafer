package v1

import (
	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/xh3b4sd/wafer/service/buyer"
	"github.com/xh3b4sd/wafer/service/buyer/runtime"
	"github.com/xh3b4sd/wafer/service/buyer/runtime/config"
	"github.com/xh3b4sd/wafer/service/buyer/runtime/state"
	"github.com/xh3b4sd/wafer/service/informer"
)

// Config is the configuration used to create a new buyer.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Runtime config.Config
}

// DefaultConfig returns the default configuration used to create a new buyer
// by best effort.
func DefaultConfig() Config {
	runtimeConfig := config.Config{}
	runtimeConfig.Trade.Concurrent = 3

	config := Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Runtime: runtimeConfig,
	}

	return config
}

// New creates a new configured buyer.
func New(config Config) (buyer.Buyer, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Logger must not be empty")
	}

	// Settings.
	err := config.Runtime.Validate()
	if err != nil {
		return nil, microerror.MaskAnyf(invalidConfigError, err.Error())
	}

	newBuyer := &Buyer{
		// Dependencies.
		logger: config.Logger,

		// Internals.
		runtime: runtime.Runtime{
			Config: config.Runtime,
			State:  state.State{},
		},
	}

	return newBuyer, nil
}

// Buyer implements buyer.Buyer.
type Buyer struct {
	// Dependencies.
	logger micrologger.Logger

	// Internals.
	runtime runtime.Runtime
}

func (b *Buyer) Buy(price informer.Price) (bool, error) {
	// Here we want to track the state of the current situation before we execute
	// the check functions. Note that the order of these functions is important.
	// because the track functions partially on each other.
	beforeTrackFuncs := []TrackFunc{
		NewSetCurrentPrice(price),
		SetMaxCorridor,
	}

	for _, t := range beforeTrackFuncs {
		r, err := t(b.runtime)
		if err != nil {
			return false, microerror.MaskAny(err)
		}
		b.runtime = r
	}

	checkFuns := []CheckFunc{
		IsAboveMaxBuys,
		IsOutsideMaxCorridor,
		IsInsideMinTradePause,
	}

	for _, c := range checkFuns {
		ok, err := c(b.runtime)
		if err != nil {
			return false, microerror.MaskAny(err)
		}
		if ok {
			return false, nil
		}
	}

	// state tracking
	b.runtime.State.Trade.Concurrent++
	b.runtime.State.Trade.Price.Last = b.runtime.State.Trade.Price.Current

	return true, nil
}

func (b *Buyer) DecrTradeConcurrent() {
	b.runtime.State.Trade.Concurrent--
}

func (b *Buyer) Runtime() runtime.Runtime {
	return b.runtime
}
