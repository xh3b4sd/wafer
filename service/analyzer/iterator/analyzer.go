package iterator

import (
	microerror "github.com/giantswarm/microkit/error"

	"github.com/xh3b4sd/wafer/service/analyzer"
	"github.com/xh3b4sd/wafer/service/decider"
)

// Config is the configuration used to create a new analyzer.
type Config struct {
	// Settings.
	Foo string
}

// DefaultConfig returns the default configuration used to create a new analyzer
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Settings.
		Foo: "",
	}
}

// New creates a new configured analyzer.
func New(config Config) (analyzer.Analyzer, error) {
	// Settings.
	if config.Foo == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Foo must not be empty")
	}

	newAnalyzer := &Analyzer{
		// Internals.
		bestConfig: decider.Config{},
	}

	return newAnalyzer, nil
}

// Analyzer implements analyzer.Analyzer.
type Analyzer struct {
	// Internals.
	bestConfig decider.Config
}

func (a *Analyzer) Adjust(configA, configB decider.Config) decider.Config {
	return decider.Config{}
}

func (a *Analyzer) Config() decider.Config {
	return a.bestConfig
}

func (a *Analyzer) Iterate(config decider.Config) (decider.Config, error) {
	return decider.Config{}, nil
}
