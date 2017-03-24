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
	// Buy returns a channel which can be used to watch for buy events. A buy
	// event indicates that the watched stock market is suitable to buy
	// commodities.
	Buy() chan informer.Price
	// Consume takes an incoming price event to analyze the stock market situation
	// to identify probabilities of buy events.
	Consume(price informer.Price) error
	// Runtime returns a copy of information about the current runtime of the
	// buyer.
	Runtime() runtime.Runtime
}
