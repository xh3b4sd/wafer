// Package informer provides the interface which has to be implemented to gather
// information about market prices. Informer may provide historical data by
// reading CSV files, or realtime data by querying external APIs.
package informer

import (
	"time"
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
	// Prices returns a channel which can be used to watch market prices.
	Prices() chan Price
}
