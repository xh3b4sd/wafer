package trade

import (
	microerror "github.com/giantswarm/microkit/error"

	"github.com/xh3b4sd/wafer/service/seller/runtime/config/trade/duration"
	"github.com/xh3b4sd/wafer/service/seller/runtime/config/trade/fee"
	"github.com/xh3b4sd/wafer/service/seller/runtime/config/trade/revenue"
)

type Trade struct {
	Duration duration.Duration `json:"duration"`
	Fee      fee.Fee           `json:"fee"`
	Revenue  revenue.Revenue   `json:"revenue"`
}

func (t Trade) Validate() error {
	var err error

	err = t.Duration.Validate()
	if err != nil {
		return microerror.MaskAnyf(invalidConfigError, err.Error())
	}
	err = t.Fee.Validate()
	if err != nil {
		return microerror.MaskAnyf(invalidConfigError, err.Error())
	}
	err = t.Revenue.Validate()
	if err != nil {
		return microerror.MaskAnyf(invalidConfigError, err.Error())
	}

	return nil
}
