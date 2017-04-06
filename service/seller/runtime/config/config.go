package config

import (
	microerror "github.com/giantswarm/microkit/error"

	"github.com/xh3b4sd/wafer/service/seller/runtime/config/trade"
)

type Config struct {
	Trade trade.Trade `json:"trade"`
}

func (c Config) Validate() error {
	err := c.Trade.Validate()
	if err != nil {
		return microerror.MaskAnyf(invalidConfigError, err.Error())
	}

	return nil
}
