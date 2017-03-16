package performer

import (
	"fmt"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/xh3b4sd/wafer/service/decider"
	"github.com/xh3b4sd/wafer/service/informer"
)

// Config is the configuration used to create a new decider.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	DeciderConfig decider.Config
	Foo           string
}

// DefaultConfig returns the default configuration used to create a new decider
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		DeciderConfig: decider.Config{},
		Foo:           "",
	}
}

// New creates a new configured decider.
func New(config Config) (decider.Decider, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

	// Settings.
	if config.Foo == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Foo must not be empty")
	}

	newAnalyzer := &Decider{
		// Internals.
		buyChan:    make(chan informer.Price, 10),
		buyEvent:   true,
		checkpoint: informer.Price{},
		config:     config.DeciderConfig,
		sellChan:   make(chan informer.Price, 10),
		window:     []informer.Price{},
	}

	return newAnalyzer, nil
}

// Decider implements decider.Decider.
type Decider struct {
	// Dependencies.
	logger micrologger.Logger

	// Internals.
	buyChan    chan informer.Price
	buyEvent   bool
	checkpoint informer.Price
	config     decider.Config
	sellChan   chan informer.Price
	window     []informer.Price
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

func (d *Decider) Watch(prices chan informer.Price) {
	emptyPrice := informer.Price{}

	for price := range prices {
		// We want to coordinate our research through the price event history using
		// checkpoints. The first checkpoint has to be defined by the first price
		// event of the given queue. In case our internal checkpoint is empty, we
		// know we receive the very first event and assign it.
		if d.checkpoint == emptyPrice {
			d.checkpoint = price
		}

		go func(price informer.Price) {
			err := d.watch(price)
			if err != nil {
				d.logger.Log("error", fmt.Sprintf("%#v", err))
			}
		}(price)
	}
}

func (d *Decider) watch(price informer.Price) error {
	var err error

	d.window = append(d.window, price)

	d.window, err = calculateWindow(d.window, d.config.Analyzer.Chart.Window)
	if err != nil {
		return microerror.MaskAny(err)
	}
	// calc chart view
	// calc chart checkpoint

	// if buy event
	//     calc chart surge
	//     if chart surge big enough
	//         buy
	//     if chart surge NOT big enough
	//         do nothing

	// if sell event
	//     calc revenue
	//     if revenue big enough
	//         sell
	//     if revenue NOT big enough
	//         do nothing

	return nil
}
