package config

import (
	"github.com/xh3b4sd/wafer/service/seller/runtime/config/trade"
)

type Config struct {
	Trade trade.Trade
}

// TODO
func (c Config) Validate() error {
	return nil
}
