package v1

import (
	"github.com/xh3b4sd/wafer/service/buyer/runtime"
	"github.com/xh3b4sd/wafer/service/informer"
)

type CheckFunc func(p informer.Price, r runtime.Runtime) (bool, error)

// IsOutsideCorridor implements CheckFunc to make sure buy events do not happen
// outside a configured price range. E.g. when the price is higher than ever
// seen, it is not likely to rise even more. Then we do not want to buy.
//
// TODO check effectiveness of corridor percentiles. For now we only have a
// static absolute number which is not very representative. Percentiles may
// perform better.
func IsOutsideCorridor(p informer.Price, r runtime.Runtime) (bool, error) {
	currentPrice := p.Buy
	maxPercAllowed := r.Config.Trade.Corridor.Max
	maxPriceSeen := r.State.Trade.Corridor.Max

	maxPriceAllowed := maxPriceSeen * maxPercAllowed / 100
	isOutsideCorridor := currentPrice > maxPriceAllowed

	return isOutsideCorridor, nil
}
