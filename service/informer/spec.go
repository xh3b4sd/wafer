// Package informer provides the interface which has to be implemented to gather
// information about market prices. Informer may provide historical data by
// reading CSV files, or realtime data by querying external APIs.
package informer

import (
	"time"

	"github.com/xh3b4sd/wafer/service/informer/csv/runtime"
)

// Price holds information about a statistical event within some market.
type Price struct {
	// Buy is the buy price at a certain time.
	Buy float64
	// Sell is the sell price at a certain time.
	Sell float64
	// Time is the time at which a certain buy and sell price occured.
	Time time.Time
}

// Informer provides market prices.
type Informer interface {
	// Prices returns a list of channels which can be used to watch market prices.
	// Each channel of the returned list represents price events from its
	// logically own chart. E.g. the CSV informer implements consuming multiple
	// CSV files. Each CSV file contains price events of potentially different
	// stock markets. Also note that the returned list of channels must be
	// consumed beginning with the first channel of the list. Consuming the last
	// channel of the returned list at first would cause a dead lock.
	Prices() []chan Price
	// Runtime returns a copy of information about the current runtime of the
	// informer.
	Runtime() runtime.Runtime
}
