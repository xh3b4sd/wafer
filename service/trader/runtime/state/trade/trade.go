package trade

import (
	"github.com/xh3b4sd/wafer/service/trader/runtime/state/trade/revenue"
)

type Trade struct {
	// Cycles is the number of buy and sell iterations. After one buy must come
	// one sell.
	Cycles  int64
	Revenue revenue.Revenue
}
