// Package analyzer implements analytical fucntionalities to optimize actual
// trade results. That is, increasing revenue.
package analyzer

import (
	"bytes"

	"github.com/xh3b4sd/wafer/service/analyzer/runtime"
)

// Analyzer provides statistical calculations to optimize configurations for
// buyers, sellers and traders.
type Analyzer interface {
	// Execute runs the exchange continuously and blocks until the configured
	// client does not provide any further price events.
	Execute()
	// Runtime returns a copy of the current statistical information about the
	// current analyzer process.
	Runtime() runtime.Runtime
}

// Visuale draws graphs of charts and their trade events.
type Visualizer interface {
	// Render returns a buffer containing the rendered data of a chart's PNG
	// image.
	Render() (*bytes.Buffer, error)
}
