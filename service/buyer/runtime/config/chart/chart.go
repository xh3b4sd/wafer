package chart

import (
	"time"
)

// Chart holds information about the chart configuration.
type Chart struct {
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
