package analyzer

import (
	"bytes"

	"github.com/xh3b4sd/wafer/service/decider"
)

// Analyzer automatically calculates the best performing decider configuration.
type Analyzer interface {
	// Adjust modifies the given decider configuration by a small degree and
	// returns the modified version.
	Adjust(config decider.Config) (decider.Config, error)
	// Config returns the decider configuration producing the best known
	// performance.
	Config() decider.Config
	// Iterate executes a decider with the given configuration and adjusts it
	// slightly according to its performance. The resulting decider configuration
	// can be used to call Iterate again to further optimize the decider
	// configuration.
	Iterate(config decider.Config) (decider.Config, error)
	// Render returns a buffer containing the rendered data of a chart's PNG
	// image.
	Render() *bytes.Buffer
	// Revenue returns the best revenue produced by the analyzer so far.
	Revenue() float64
}
