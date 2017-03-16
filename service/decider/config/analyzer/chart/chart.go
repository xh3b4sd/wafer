package chart

import (
	"time"
)

// Chart holds information about the analyzer chart configuration. The chart
// window configuration can be visualized as described below to help understand
// the several settings.
//
//                                  <------------ Chart.View ------------>
//         <-- Chart.Convolution -->
//         <------------ Chart.View ------------>
//         <----------------------- Chart.Window ----------------------->
//
//    <----------------------- complete chart history ----------------------->
//
type Chart struct {
	// Convolution is the sliding step within the chart window used to determine
	// the next chart view.
	Convolution time.Duration
	// View is a single view being analyzed.
	View time.Duration
	// Window is the complete buffered chart history separated into single views.
	Window time.Duration
}
