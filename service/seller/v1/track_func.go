package v1

import (
	"github.com/xh3b4sd/wafer/service/informer"
	"github.com/xh3b4sd/wafer/service/seller/runtime"
)

type TrackFunc func(r runtime.Runtime) (runtime.Runtime, error)

// NewSetCurrentDuration returns a new function which implements TrackFunc to
// set the current duration to the runtime state.
func NewSetCurrentDuration(currentPrice, buyPrice informer.Price) TrackFunc {
	return func(r runtime.Runtime) (runtime.Runtime, error) {
		r.State.Trade.Duration = currentPrice.Time.Sub(buyPrice.Time)

		return r, nil
	}
}

// NewSetCurrentRevenue returns a new function which implements TrackFunc to set
// the current revenue to the runtime state.
func NewSetCurrentRevenue(currentPrice, buyPrice informer.Price) TrackFunc {
	return func(r runtime.Runtime) (runtime.Runtime, error) {
		r.State.Trade.Revenue = calculateRevenue(buyPrice.Buy, currentPrice.Sell, r.Config.Trade.Fee.Min)

		return r, nil
	}
}
