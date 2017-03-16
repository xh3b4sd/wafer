package trader

import (
	"github.com/xh3b4sd/wafer/service/decider/config/trader/duration"
	"github.com/xh3b4sd/wafer/service/decider/config/trader/fee"
	"github.com/xh3b4sd/wafer/service/decider/config/trader/revenue"
)

type Trader struct {
	Duration duration.Duration
	Fee      fee.Fee
	Revenue  revenue.Revenue
}
