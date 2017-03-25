package config

import (
	"github.com/xh3b4sd/wafer/service/buyer/runtime/config/chart"
	"github.com/xh3b4sd/wafer/service/buyer/runtime/config/surge"
	"github.com/xh3b4sd/wafer/service/buyer/runtime/config/trade"
)

type Config struct {
	Chart chart.Chart
	Surge surge.Surge
	Trade trade.Trade
}
