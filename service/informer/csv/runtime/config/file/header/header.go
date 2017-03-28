package header

import (
	microerror "github.com/giantswarm/microkit/error"
)

type Header struct {
	// Buy is the index of the row representing buy prices within the given
	// CSV file.
	Buy int
	// Ignore decides whether to ignore the first line of the given CSV.
	// This can be set to true in case the first line does not represent actual
	// data.
	Ignore bool
	// Sell is the index of the row representing sell prices within the given
	// CSV file.
	Sell int
	// Time is the index of the row representing price times within the given
	// CSV file.
	Time int
}

func (h Header) Validate() error {
	if h.Buy == h.Sell {
		return microerror.MaskAnyf(invalidConfigError, "h.Buy must not be equal to h.Sell")
	}
	if h.Buy == h.Time {
		return microerror.MaskAnyf(invalidConfigError, "h.Buy must not be equal to h.Time")
	}
	if h.Sell == h.Time {
		return microerror.MaskAnyf(invalidConfigError, "h.Sell must not be equal to h.Time")
	}

	return nil
}
