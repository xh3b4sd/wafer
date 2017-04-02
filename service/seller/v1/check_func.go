package v1

import (
	"github.com/xh3b4sd/wafer/service/seller/runtime"
)

type CheckFunc func(r runtime.Runtime) (bool, error)

// IsBelowMinTradeDuration implements CheckFunc to make sure sell events do not
// happen under a configured trade duration. E.g. when the buy price event is
// not long enough ago, we do not want to sell.
func IsBelowMinTradeDuration(r runtime.Runtime) (bool, error) {
	isBelowMinTradeDuration := r.State.Trade.Duration < r.Config.Trade.Duration.Min

	return isBelowMinTradeDuration, nil
}

// IsBelowMinTradeRevenue implements CheckFunc to make sure sell events do not
// happen under a configured trade revenue. E.g. when the sell price lower than
// the buy price, we do not want to sell.
func IsBelowMinTradeRevenue(r runtime.Runtime) (bool, error) {
	isBelowMinTradeRevenue := r.State.Trade.Revenue < r.Config.Trade.Revenue.Min

	return isBelowMinTradeRevenue, nil
}
