// Package trader provides the interface which has to be implemented to combine
// the business logic of buyers and sellers. This can enable the participation
// on some market to increase revenue. Trader might buy and sell bitcoins,
// ether, or other conservative stock market shares. Trader may be used to
// operate against analyzation exchanges to optimize certain algorithm
// configurations.
package trader

import (
	"github.com/xh3b4sd/wafer/service/trader/runtime"
)

// Trader combines the business logic of informers, buyers, sellers and clients.
// All kinds of informers can be used to feed the buyers and sellers. That way
// some CSV data can be consumed or some realtime events from some real stock
// exchange API. All kinds of buyers and sellers can be used which provide all
// kinds of ruler functions. All kinds of clients can be used to forward buy and
// sell events to all kinds of exchanges. This can be used to feed some local
// analyzer exchange or some real stock exchange API.
type Trader interface {
	// Execute runs the trader continuously and blocks until the configured
	// informer does not provide any further price events.
	Execute() error
	// Runtime returns a copy of the current statistical information about the
	// current trader process.
	Runtime() runtime.Runtime
}
