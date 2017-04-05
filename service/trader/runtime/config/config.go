package config

import (
	"github.com/xh3b4sd/wafer/service/trader/runtime/config/trade"
)

type Config struct {
	Trade trade.Trade `json:"trader"`
}

// TODO
func (c Config) Validate() error {
	return nil
}
