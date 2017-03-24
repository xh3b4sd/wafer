package state

import (
	"github.com/xh3b4sd/wafer/service/seller/runtime/state/chart"
	"github.com/xh3b4sd/wafer/service/seller/runtime/state/trade"
)

type State struct {
	Chart chart.Chart
	Trade trade.Trade
}
