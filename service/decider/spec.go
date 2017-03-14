// Package decider provide the interface which has to be implemented to judge
// situation on some market. Decider should be performance based. Decider may
// leverage neural networks or more static algorithms. Deciders should compete
// with each other.
package decider

import (
	"time"
)

// Config holds information about decider settings used to judge some market
// situation.
type Config struct {
	// MaxDuration is the maximum time a single trade is allowed to take.
	MaxDuration time.Duration
	// MinDuration is the minimum time a single trade is allowed to take.
	MinDuration time.Duration
	// MaxRevenue is the maximum revenue a single trade is allowed to make.
	MaxRevenue float64
	// MinRevenue is the minimum revenue a single trade is allowed to make.
	MinRevenue float64
}

// Decider judges based on events to qualify if the situation at some market is
// suited to either buy or sell commodities.
type Decider interface {
	// Buy returns a channel which can be used to watch for buy events. A buy
	// event indicates that the market is suitable to buy commodities. Once there
	// was a buy event, a sell event must follow at some point. Two buy events
	// without a sell event in between must never happen.
	Buy() chan struct{}
	// Config returns a copy of information about the current settings of the
	// decider. The settings represented here must never change over the lifetime
	// of a decider.
	Config() Config
	// Sell returns a channel which can be used to watch for sell events. A sell
	// event indicates that the market is suitable to sell commodities. Once there
	// was a sell event, a buy event must follow at some point. Two sell events
	// without a buy event in between must never happen.
	Sell() chan struct{}
}
