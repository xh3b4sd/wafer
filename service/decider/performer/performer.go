package performer

import (
	microerror "github.com/giantswarm/microkit/error"

	"github.com/xh3b4sd/wafer/service/decider"
	"github.com/xh3b4sd/wafer/service/informer"
)

// Config is the configuration used to create a new decider.
type Config struct {
	// Settings.
	DeciderConfig decider.Config
	Foo           string
}

// DefaultConfig returns the default configuration used to create a new decider
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Settings.
		DeciderConfig: decider.Config{},
		Foo:           "",
	}
}

// New creates a new configured decider.
func New(config Config) (decider.Decider, error) {
	// Settings.
	if config.Foo == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Foo must not be empty")
	}

	newAnalyzer := &Decider{
		// Internals.
		buyChan:  make(chan informer.Price, 10),
		config:   config.DeciderConfig,
		sellChan: make(chan informer.Price, 10),
	}

	return newAnalyzer, nil
}

// Decider implements decider.Decider.
type Decider struct {
	// Internals.
	buyChan  chan informer.Price
	config   decider.Config
	sellChan chan informer.Price
}

func (d *Decider) Buy() chan informer.Price {
	return d.buyChan
}

func (d *Decider) Config() decider.Config {
	return d.config
}

func (d *Decider) Sell() chan informer.Price {
	return d.sellChan
}
