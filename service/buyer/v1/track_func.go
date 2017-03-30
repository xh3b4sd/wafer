package v1

import (
	microerror "github.com/giantswarm/microkit/error"

	"github.com/xh3b4sd/wafer/service/buyer/runtime"
	"github.com/xh3b4sd/wafer/service/informer"
)

type TrackFunc func(r runtime.Runtime) (runtime.Runtime, error)

// NewSetCurrentPrice returns a new function which implements TrackFunc to set
// the current price to the runtime state.
func NewSetCurrentPrice(p informer.Price) TrackFunc {
	return func(r runtime.Runtime) (runtime.Runtime, error) {
		r.State.Trade.Price = p

		return r, nil
	}
}

func SetChartWindow(r runtime.Runtime) (runtime.Runtime, error) {
	var err error
	r.State.Chart.Window = append(r.State.Chart.Window, r.State.Trade.Price)
	r.State.Chart.Window, err = calculateWindow(r.State.Chart.Window, r.Config.Chart.Window)
	if err != nil {
		return runtime.Runtime{}, microerror.MaskAny(err)
	}

	return r, nil
}

func SetLastSurge(r runtime.Runtime) (runtime.Runtime, error) {
	r.State.Surge.Last = findLastSurge(r.State.Chart.Window, r.Config.Surge.Tolerance)

	return r, nil
}

func SetMaxCorridor(r runtime.Runtime) (runtime.Runtime, error) {
	if r.State.Trade.Corridor.Max < r.State.Trade.Price.Buy {
		r.State.Trade.Corridor.Max = r.State.Trade.Price.Buy
	}

	return r, nil
}
