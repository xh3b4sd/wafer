package v1

import (
	"github.com/xh3b4sd/wafer/service/buyer/runtime"
)

type CheckFunc func(r runtime.Runtime) (bool, error)

// IsAboveMaxBuys implements CheckFunc to make sure too many buy events do not
// happen at the same time. E.g. we cannot run out of budget because we buy too
// many times without selling before to free invested budget.
func IsAboveMaxBuys(r runtime.Runtime) (bool, error) {
	isAboveMaxBuys := r.State.Trade.Concurrent >= r.Config.Trade.Concurrent

	return isAboveMaxBuys, nil
}

// IsOutsideMaxCorridor implements CheckFunc to make sure buy events do not
// happen outside a configured price range. E.g. when the price is higher than
// ever seen, it is not likely to rise even more. Then we do not want to buy.
//
// TODO check effectiveness of corridor percentiles. For now we only have a
// static absolute number which is not very representative. Percentiles may
// perform better.
func IsOutsideMaxCorridor(r runtime.Runtime) (bool, error) {
	currentPrice := r.State.Trade.Price.Current.Buy
	maxPercAllowed := r.Config.Trade.Corridor.Max
	maxPriceSeen := r.State.Trade.Corridor.Max

	maxPriceAllowed := maxPriceSeen * maxPercAllowed / 100
	isAboveMaxTradeLimit := currentPrice > maxPriceAllowed

	return isAboveMaxTradeLimit, nil
}

// IsInsideMinTradePause implements CheckFunc to make sure buy events do not
// happen under a configured trade pause. E.g. we want to wait some time after
// buying commodities before we buy again.
func IsInsideMinTradePause(r runtime.Runtime) (bool, error) {
	if r.State.Trade.Price.Last.Time.IsZero() {
		return false, nil
	}

	isInsideMinTradePause := r.State.Trade.Price.Last.Time.Add(r.Config.Trade.Pause.Min).After(r.State.Trade.Price.Current.Time)

	return isInsideMinTradePause, nil
}
