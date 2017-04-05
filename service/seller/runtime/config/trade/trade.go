package trade

import (
	"github.com/xh3b4sd/wafer/service/seller/runtime/config/trade/duration"
	"github.com/xh3b4sd/wafer/service/seller/runtime/config/trade/fee"
	"github.com/xh3b4sd/wafer/service/seller/runtime/config/trade/revenue"
)

type Trade struct {
	Duration duration.Duration `json:"duration"`
	Fee      fee.Fee           `json:"fee"`
	Revenue  revenue.Revenue   `json:"revenue"`
}
