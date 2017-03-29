package v1

import (
	"time"

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
	// state tracking
	var err error
	b.runtime.State.Chart.Window = append(b.runtime.State.Chart.Window, price)
	b.runtime.State.Chart.Window, err = calculateWindow(b.runtime.State.Chart.Window, b.runtime.Config.Chart.Window)
	if err != nil {
		return false, microerror.MaskAny(err)
	}

	prices := findLastSurge(b.runtime.State.Chart.Window, b.runtime.Config.Surge.Tolerance)

	// state tracking
	if b.runtime.State.Trade.Corridor.Max < price.Buy {
		b.runtime.State.Trade.Corridor.Max = price.Buy
	}

	// rule checking
	isOutsideCorridor, err := IsOutsideCorridor(price, b.runtime)
	if err != nil {
		return false, microerror.MaskAny(err)
	}
	if isOutsideCorridor {
		return false, nil
	}

	// rule checking
	surge := calculateSurgeAverage(prices)
	if surge < b.runtime.Config.Surge.Min {
		return false, nil
	}

	// rule checking
	duration := calculateSurgeDuration(prices)
	duration = duration + (time.Duration(duration.Seconds()*surge*surge/100) * time.Second)
	if duration < b.runtime.Config.Surge.Duration.Min {
		return false, nil
	}

	// rule checking
	if !b.runtime.State.Trade.Buy.Last.IsZero() && b.runtime.State.Trade.Buy.Last.Add(b.runtime.Config.Trade.Pause.Min).Before(price.Time) {
		return false, nil
	}

	// state tracking
	b.runtime.State.Trade.Buy.Last = price.Time

	return true, nil
}

func (b *Buyer) Runtime() runtime.Runtime {
	return b.runtime
}
