package trade

import (
	"github.com/xh3b4sd/wafer/service/buyer/runtime/config/trade/corridor"
	"github.com/xh3b4sd/wafer/service/buyer/runtime/config/trade/pause"
)

type Trade struct {
	Corridor corridor.Corridor
	Pause    pause.Pause
}
