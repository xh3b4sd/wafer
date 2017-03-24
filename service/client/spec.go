package client

import (
	"github.com/xh3b4sd/wafer/service/informer"
)

// Client communicates with some exchange to actually process a buy or a sell.
type Client interface {
	// Buy notifies the configured exchange to buy for a certain price.
	Buy(price informer.Price) error
	// Close shuts down the client and all of its activity.
	Close() error
	// Sell notifies the configured exchange to sell for a certain price.
	Sell(price informer.Price) error
}
