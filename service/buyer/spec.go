// Package buyer provides the interface which has to be implemented to judge
// buying situations on some stock market. The buyer may leverage neural
// networks or more static algorithms. Buyers should compete with each other.
package buyer

import (
	"github.com/xh3b4sd/wafer/service/buyer/runtime"
	"github.com/xh3b4sd/wafer/service/informer"
)

// Buyer judges based on events to qualify if the situation at some stock market
// is suited to buy commodities.
type Buyer interface {
	// Consume takes an incoming price event to analyze the stock market situation
	// to identify probabilities of buy events. In case Sell returns true, a sell
	// event is intended to happen. A buy event indicates that the watched stock
	// market is suitable to buy commodities.
	Buy(price informer.Price) (bool, error)
	// DecrTradeConcurrent is to make the buyer aware of an sell event which implies
	// that the buyer must decrement the internal state counter of the number of
	// concurrent buy events.
	DecrTradeConcurrent()
	// Runtime returns a copy of information about the current runtime of the
	// buyer.
	Runtime() runtime.Runtime
}
