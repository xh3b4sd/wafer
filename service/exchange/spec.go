// Package exchange ...
package exchange

import (
	"bytes"

	"github.com/xh3b4sd/wafer/service/informer"
)

// Exchange is fed with price events which are provided by a configured client.
// An exchange can have very different purposes. A possible implementation may
// provide statistical calculations to optimize configurations for buyers and
// sellers.
type Exchange interface {
	// Buys returns a list of buy price events.
	Buys() []informer.Price
	// Close shuts down the exchange.
	Close() error
	// Execute runs the exchange continuously and blocks until the configured
	// client does not provide any further price events.
	Execute()
	// Render returns a buffer containing the rendered data of a chart's PNG
	// image.
	Render() (*bytes.Buffer, error)
	// Sells returns a list of sell price events.
	Sells() []informer.Price
}
