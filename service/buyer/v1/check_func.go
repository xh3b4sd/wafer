package v1

import (
	"github.com/xh3b4sd/wafer/service/buyer/runtime"
)

type CheckFunc func(r runtime.Runtime) (bool, error)

// IsAboveMaxTradeLimit implements CheckFunc to make sure buy events do not
// happen outside a configured price range. E.g. when the price is higher than
// ever seen, it is not likely to rise even more. Then we do not want to buy.
//
// TODO check effectiveness of corridor percentiles. For now we only have a
// static absolute number which is not very representative. Percentiles may
// perform better.
func IsAboveMaxTradeLimit(r runtime.Runtime) (bool, error) {
	currentPrice := r.State.Trade.Price.Buy
	maxPercAllowed := r.Config.Trade.Corridor.Max
	maxPriceSeen := r.State.Trade.Corridor.Max

	maxPriceAllowed := maxPriceSeen * maxPercAllowed / 100
	isAboveMaxTradeLimit := currentPrice > maxPriceAllowed

	return isAboveMaxTradeLimit, nil
}

// IsUnderMinSurgeAngle implements CheckFunc to make sure buy events do not
// happen under a configured surge angle. E.g. when the price is not rising, we
// do not want to buy.
func IsUnderMinSurgeAngle(r runtime.Runtime) (bool, error) {
	surge := calculateSurgeAverage(r.State.Surge.Last)
	isUnderMinSurgeAngle := surge < r.Config.Surge.Min

	return isUnderMinSurgeAngle, nil
}

// IsUnderMinSurgeDuration implements CheckFunc to make sure buy events do not
// happen under a configured surge duration. E.g. when the price is not rising
// for a certain time range, we do not want to buy.
func IsUnderMinSurgeDuration(r runtime.Runtime) (bool, error) {
	duration := calculateSurgeDuration(r.State.Surge.Last)
	isUnderMinSurgeDuration := duration < r.Config.Surge.Duration.Min

	return isUnderMinSurgeDuration, nil
}

// IsUnderMinTradePause implements CheckFunc to make sure buy events do not
// happen under a configured trade pause. E.g. we want to wait some time after
// buying commodities before we buy again.
func IsUnderMinTradePause(r runtime.Runtime) (bool, error) {
	if r.State.Trade.Buy.Last.Time.IsZero() {
		return false, nil
	}

	isUnderMinTradePause := r.State.Trade.Buy.Last.Time.Add(r.Config.Trade.Pause.Min).Before(r.State.Trade.Price.Time)

	return isUnderMinTradePause, nil
}
