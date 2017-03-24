// Package seller provides the interface which has to be implemented to judge
// selling situations on some stock market. The seller may leverage neural
// networks or more static algorithms. Sellers should compete with each other.
package seller

import (
	"github.com/xh3b4sd/wafer/service/informer"
	"github.com/xh3b4sd/wafer/service/seller/runtime"
)

// Seller judges based on events to qualify if the situation at some stock
// market is suited to sell commodities.
type Seller interface {
	// Runtime returns a copy of information about the current runtime of the
	// seller.
	Runtime() runtime.Runtime
	// Sell takes the actual buy price and the currently incoming price event to
	// analyze the stock market situation to identify probabilities of sell
	// events. In case Sell returns true, a sell event is intended to happen. A
	// sell event indicates that the watched stock market is suitable to sell
	// commodities.
	Sell(buyPrice, currentPrice informer.Price) (bool, error)
}
