package trade

import (
	microerror "github.com/giantswarm/microkit/error"

	"github.com/xh3b4sd/wafer/service/buyer/runtime/config/trade/corridor"
	"github.com/xh3b4sd/wafer/service/buyer/runtime/config/trade/pause"
)

type Trade struct {
	// Concurrent is the maximum number of allowed parallel buy events.
	Concurrent int               `json:"concurrent"`
	Corridor   corridor.Corridor `json:"corridor"`
	Pause      pause.Pause       `json:"pause"`
}

func (t Trade) Validate() error {
	if t.Concurrent == 0 {
		return microerror.MaskAnyf(invalidConfigError, "Trade.Concurrent must not be empty")
	}

	var err error

	err = t.Corridor.Validate()
	if err != nil {
		return microerror.MaskAnyf(invalidConfigError, err.Error())
	}
	err = t.Pause.Validate()
	if err != nil {
		return microerror.MaskAnyf(invalidConfigError, err.Error())
	}

	return nil
}
