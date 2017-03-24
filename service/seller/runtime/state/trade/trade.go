package trade

import (
	"github.com/xh3b4sd/wafer/service/informer"
	"github.com/xh3b4sd/wafer/service/seller/runtime/state/trade/duration"
	"github.com/xh3b4sd/wafer/service/seller/runtime/state/trade/revenue"
)

type Trade struct {
	// Buy is the buy price the seller takes for granted to calculate
	// probabilities for sell events.
	Buy      informer.Price
	Duration duration.Duration
	Revenue  revenue.Revenue
	// Total is the total number of trades finsihed by the seller.
	Total int
}
