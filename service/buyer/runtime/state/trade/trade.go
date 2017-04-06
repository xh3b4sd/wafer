package trade

import (
	"github.com/xh3b4sd/wafer/service/buyer/runtime/state/trade/corridor"
	"github.com/xh3b4sd/wafer/service/buyer/runtime/state/trade/price"
)

type Trade struct {
	Corridor corridor.Corridor
	Price    price.Price
	// Concurrent is the number of concurrent buys emitted by the buyer.
	Concurrent int
}
