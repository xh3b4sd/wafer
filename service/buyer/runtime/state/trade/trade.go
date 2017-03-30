package trade

import (
	"github.com/xh3b4sd/wafer/service/buyer/runtime/state/trade/buy"
	"github.com/xh3b4sd/wafer/service/buyer/runtime/state/trade/corridor"
	"github.com/xh3b4sd/wafer/service/informer"
)

type Trade struct {
	Buy      buy.Buy
	Corridor corridor.Corridor
	// Price is the current price event which can be used to buy commodities.
	Price informer.Price
}
