package chart

import (
	"github.com/xh3b4sd/wafer/service/informer"
)

// Chart holds information about the chart configuration.
type Chart struct {
	// Window holds all price events within the time span ranging from the
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
	Window []informer.Price
}
