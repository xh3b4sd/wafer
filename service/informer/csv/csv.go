// Package csv provides the implementation of an informer able to read CSV
// files.
package csv

import (
	microerror "github.com/giantswarm/microkit/error"

	"github.com/xh3b4sd/wafer/service/informer"
	"github.com/xh3b4sd/wafer/service/informer/csv/runtime"
	runtimeconfigdir "github.com/xh3b4sd/wafer/service/informer/csv/runtime/config/dir"
	runtimeconfigfile "github.com/xh3b4sd/wafer/service/informer/csv/runtime/config/file"
	runtimestatefile "github.com/xh3b4sd/wafer/service/informer/csv/runtime/state/file"
	stateprice "github.com/xh3b4sd/wafer/service/informer/csv/runtime/state/price"
)

// Config is the configuration used to create a new informer.
type Config struct {
	// Settings.

	// Dir is the config for an absolute location of the CSV dir to consume.
	// Either File or dir can be used at the same time.
	Dir runtimeconfigdir.Dir
	// File is the config for an absolute location of the CSV file to consume.
	// Either File or dir can be used at the same time.
	File runtimeconfigfile.File
}

// DefaultConfig returns the default configuration used to create a new informer
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Settings.
		Dir:  runtimeconfigdir.Dir{},
		File: runtimeconfigfile.File{},
	}
}

// New creates a new configured informer.
func New(config Config) (informer.Informer, error) {
	// Settings.
	dErr := config.Dir.Validate()
	fErr := config.File.Validate()
	if dErr != nil && fErr != nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "either config.Dir or config.File must be given")
	}
	if dErr == nil && fErr == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "either config.Dir or config.File must be given")
	}

	var files []runtimestatefile.File
	var err error

	if dErr != nil {
		files, err = fileToFiles(config.File)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	} else {
		files, err = dirToFiles(config.Dir)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	prices, err := filesToPrices(files)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	newInformer := &Informer{
		// Internals.
		prices:  prices,
		runtime: runtime.Runtime{},
	}

	for _, p := range prices {
		if len(p) < 2 {
			return nil, microerror.MaskAnyf(invalidConfigError, "chart must contain at least 2 price events")
		}

		price := stateprice.Price{
			End:    p[len(p)-1].Time,
			Events: len(p),
			Start:  p[0].Time,
		}
		newInformer.runtime.State.Prices = append(newInformer.runtime.State.Prices, price)
	}

	return newInformer, nil
}

// Informer implements informer.Informer.
type Informer struct {
	// Internals.
	prices  [][]informer.Price
	runtime runtime.Runtime
}

// Prices returns a list of price channels containing price events. These hold
// buy and sell prices as well as their corresponding timestamps. Note that buy
// and sell prices are parsed as float64 and the CSV informer assumes the
// timestamp is a usual unix timestamp in seconds. Also note that the returned
// list of channels must be consumed beginning with the first channel of the
// list. Consuming the last channel of the returned list at first would cause a
// dead lock.
func (i *Informer) Prices() []chan informer.Price {
	prices := make([]chan informer.Price, len(i.prices))
	for i, _ := range prices {
		prices[i] = make(chan informer.Price, 10)
	}

	go func() {
		for j, c := range prices {
			for _, p := range i.prices[j] {
				c <- p
			}
			close(c)
		}
	}()

	return prices
}

func (i *Informer) Runtime() runtime.Runtime {
	return i.runtime
}
