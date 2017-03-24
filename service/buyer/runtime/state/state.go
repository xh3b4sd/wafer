package state

import (
	"github.com/xh3b4sd/wafer/service/buyer/runtime/state/chart"
	"github.com/xh3b4sd/wafer/service/buyer/runtime/state/trade"
)

type State struct {
	Chart chart.Chart
	Trade trade.Trade
}
