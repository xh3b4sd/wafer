package chart

import (
	"time"
)

// Chart holds information about the chart configuration.
type Chart struct {
	// TODO
	// Max is the biggest price seen on the watched stock market chart.
	// Max informer.Price
	// TODO
	// Min is the lowest price seen on the watched stock market chart.
	// Min informer.Price

	// Window describes all price events within the time span ranging from the
	// present back to some point in the past.
	//
	//            past                                        present
	//              |                                            |
	//              v                                            v
	//
	//              <--------------- Chart.Window --------------->
	//
	//    <--------------- complete chart history --------------->
	//
	Window time.Duration
}
