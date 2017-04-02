package trade

import (
	"time"

	"github.com/xh3b4sd/wafer/service/informer"
)

type Trade struct {
	// Buy is the buy price the seller takes for granted to calculate
	// probabilities for sell events.
	Buy informer.Price
	// Duration is the minimum time it took for the seller to emit a sell event.
	Duration time.Duration
	// Revenue is the total amount of revenue the seller made so far.
	Revenue float64
	// Total is the total number of trades finsihed by the seller.
	Total int
}
