package price

import (
	"github.com/xh3b4sd/wafer/service/informer"
)

type Price struct {
	// Current is the current price event which can be used to buy commodities.
	Current informer.Price
	// Last is the last price event which was used to buy commodities.
	Last informer.Price
}
