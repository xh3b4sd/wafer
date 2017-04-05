package config

import (
	"github.com/xh3b4sd/wafer/service/buyer/runtime/config/chart"
	"github.com/xh3b4sd/wafer/service/buyer/runtime/config/surge"
	"github.com/xh3b4sd/wafer/service/buyer/runtime/config/trade"
)

type Config struct {
	Chart chart.Chart `json:"chart"`
	Surge surge.Surge `json:"surge"`
	Trade trade.Trade `json:"trade"`
}

// TODO
func (c Config) Validate() error {
	return nil
}
