package client

import (
	"github.com/xh3b4sd/wafer/service/informer"
)

// Client communicates with some exchange to actually process a buy or a sell.
type Client interface {
	// Buy notifies the configured exchange to buy for a certain price with
	// respect to the given volume.
	Buy(price informer.Price, volume float64) error
	// Close shuts down the client and all of its activity.
	Close() error
	// Sell notifies the configured exchange to sell for a certain price with
	// respect to the given volume.
	Sell(price informer.Price, volume float64) error
}
