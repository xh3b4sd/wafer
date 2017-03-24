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
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Runtime: config.Config{},
	}
}

// New creates a new configured buyer.
func New(config Config) (buyer.Buyer, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Logger must not be empty")
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
	var err error
	b.runtime.State.Chart.Window = append(b.runtime.State.Chart.Window, price)
	b.runtime.State.Chart.Window, err = calculateWindow(b.runtime.State.Chart.Window, b.runtime.Config.Chart.Window)
	if IsNotEnoughData(err) {
		// In case there is not enough data yet, we cannot continue with the chart
		// analyzation. So we return here and wait for the next events and proceed
		// later, as soon as there is enough data for our algorithm.
		return false, nil
	} else if err != nil {
		return false, microerror.MaskAny(err)
	}

	prices := findLastSurge(b.runtime.State.Chart.Window)
	surge := calculateSurgeAverage(prices)

	// TODO find out why surge is 2.5 and not 45
	if surge < b.runtime.Config.Surge.Min {
		return false, nil
	}

	return true, nil
}

func (b *Buyer) Runtime() runtime.Runtime {
	return b.runtime
}
