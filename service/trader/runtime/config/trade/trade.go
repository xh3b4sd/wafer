package trade

import (
	microerror "github.com/giantswarm/microkit/error"
)

type Trade struct {
	// Budget is the amount of money used to buy and sell commodities.
	Budget float64 `json:"budget"`
}

func (t Trade) Validate() error {
	if t.Budget == 0 {
		return microerror.MaskAnyf(invalidConfigError, "Trade.Budget must not be empty")
	}

	return nil
}
