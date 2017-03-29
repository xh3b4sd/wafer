package trade

import (
	"github.com/xh3b4sd/wafer/service/buyer/runtime/state/trade/buy"
	"github.com/xh3b4sd/wafer/service/buyer/runtime/state/trade/corridor"
)

type Trade struct {
	Buy      buy.Buy
	Corridor corridor.Corridor
}
