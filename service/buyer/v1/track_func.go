package v1

import (
	"github.com/xh3b4sd/wafer/service/buyer/runtime"
	"github.com/xh3b4sd/wafer/service/informer"
)

type TrackFunc func(r runtime.Runtime) (runtime.Runtime, error)

// NewSetCurrentPrice returns a new function which implements TrackFunc to set
// the current price to the runtime state.
func NewSetCurrentPrice(p informer.Price) TrackFunc {
	return func(r runtime.Runtime) (runtime.Runtime, error) {
		r.State.Trade.Price.Current = p

		return r, nil
	}
}

func SetMaxCorridor(r runtime.Runtime) (runtime.Runtime, error) {
	if r.State.Trade.Corridor.Max < r.State.Trade.Price.Current.Buy {
		r.State.Trade.Corridor.Max = r.State.Trade.Price.Current.Buy
	}

	return r, nil
}
