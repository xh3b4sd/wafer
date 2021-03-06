package fee

import (
	microerror "github.com/giantswarm/microkit/error"
)

type Fee struct {
	// Min is the minimum amount of fees that have to be respected when
	// calculating the revenue of a single trade. Note that this probably has to
	// be respected for buy and sell events, which means that this number
	// represents the fee payed twice. E.g. when the fee for buy and sell events
	// is 1.49% each, the fee the seller must respect is 2.98%.
	Min float64 `json:"min"`
}

func (f Fee) Validate() error {
	if f.Min == 0 {
		return microerror.MaskAnyf(invalidConfigError, "Fee.Min must not be empty")
	}

	return nil
}
