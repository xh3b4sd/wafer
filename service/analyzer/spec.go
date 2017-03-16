package analyzer

import (
	"github.com/xh3b4sd/wafer/service/decider"
)

// Analyzer automatically calculates the best performing decider configuration.
type Analyzer interface {
	// Adjust modifies the given decider configuration by a small degree and
	// returns the modified version. Adjust is given two consecutive decider
	// configurations as calculated via Iterate. Therefore Adjust must modify the
	// second decider configuration according to the first one.
	Adjust(configA, configB decider.Config) (decider.Config, error)
	// Config returns the decider configuration producing the best known
	// performance.
	Config() decider.Config
	// Iterate executes a decider with the given configuration and adjusts it
	// slightly according to its performance. The resulting decider configuration
	// can be used to call Iterate again to further optimize the decider
	// configuration.
	Iterate(config decider.Config) (decider.Config, error)
}
