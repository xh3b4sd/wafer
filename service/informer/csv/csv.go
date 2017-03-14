// Package csv provides the implementation of an informer able to read CSV
// files.
package csv

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	microerror "github.com/giantswarm/microkit/error"

	"github.com/xh3b4sd/wafer/service/informer"
)

// Config is the configuration used to create a new informer.
type Config struct {
	// Settings.

	// BuyIndex is the index of the row representing buy prices within the given
	// CSV file.
	BuyIndex int
	// File is the absolute location of the CSV file to consume.
	File string
	// IgnoreHeader decides whether to ignore the first line of the given CSV.
	// This can be set to true in case the first line does not represent actual
	// data.
	IgnoreHeader bool
	// SellIndex is the index of the row representing sell prices within the given
	// CSV file.
	SellIndex int
	// TimeIndex is the index of the row representing price times within the given
	// CSV file.
	TimeIndex int
}

// DefaultConfig returns the default configuration used to create a new informer
// by best effort.
func DefaultConfig() Config {
	return Config{
		BuyIndex:     0,
		File:         "",
		IgnoreHeader: false,
		SellIndex:    0,
		TimeIndex:    0,
	}
}

// New creates a new configured informer.
func New(config Config) (informer.Informer, error) {
	// Settings.
	if config.BuyIndex == config.SellIndex {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.BuyIndex must not be equal to config.SellIndex")
	}
	if config.BuyIndex == config.TimeIndex {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.BuyIndex must not be equal to config.TimeIndex")
	}
	if config.SellIndex == config.TimeIndex {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.SellIndex must not be equal to config.TimeIndex")
	}
	if config.File == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.File must not be empty")
	}

	prices := make(chan informer.Price, 10)

	// Read the CSV file.
	var fields [][]string
	{
		csvFile, err := os.Open(config.File)
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
		defer csvFile.Close()

		reader := csv.NewReader(csvFile)
		fields, err = reader.ReadAll()
		if err != nil {
			return nil, microerror.MaskAny(err)
		}
	}

	// Validate the given configuration.
	{
		if config.BuyIndex >= len(fields) {
			return nil, microerror.MaskAnyf(invalidConfigError, "config.BuyIndex out of range")
		}
		if config.SellIndex >= len(fields) {
			return nil, microerror.MaskAnyf(invalidConfigError, "config.SellIndex out of range")
		}
		if config.TimeIndex >= len(fields) {
			return nil, microerror.MaskAnyf(invalidConfigError, "config.TimeIndex out of range")
		}
	}

	// Fill the prices channel.
	{
		if config.IgnoreHeader {
			fields = fields[1:]
		}

		go func() {
			for _, fs := range fields {
				b, err := strconv.ParseFloat(fs[config.BuyIndex], 64)
				if err != nil {
					panic(err)
				}
				s, err := strconv.ParseFloat(fs[config.SellIndex], 64)
				if err != nil {
					panic(err)
				}
				t, err := strconv.ParseInt(fs[config.TimeIndex], 10, 64)
				if err != nil {
					panic(err)
				}

				price := informer.Price{
					Buy:  b,
					Sell: s,
					Time: time.Unix(t, 0),
				}
				prices <- price
			}

			close(prices)
		}()
	}

	newInformer := &Informer{
		// Internals.
		prices: prices,
	}

	return newInformer, nil
}

// server manages the transport logic and endpoint registration.
type Informer struct {
	// Internals.
	prices chan informer.Price
}

// Prices returns a price channel with price events holding buy and sell prices
// as well as their corresponding timestamps. Note that buy and sell prices are
// parsed as float64 and the CSV informer assumes the timestamp is a usual unix
// timestamp in seconds.
func (i *Informer) Prices() chan informer.Price {
	return i.prices
}
