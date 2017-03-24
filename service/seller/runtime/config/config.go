package config

import (
	"github.com/xh3b4sd/wafer/service/seller/runtime/config/chart"
	"github.com/xh3b4sd/wafer/service/seller/runtime/config/decline"
	"github.com/xh3b4sd/wafer/service/seller/runtime/config/trade"
)

type Config struct {
	Chart   chart.Chart
	Decline decline.Decline
	Trade   trade.Trade
}
