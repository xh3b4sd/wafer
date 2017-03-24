package config

import (
	"github.com/xh3b4sd/wafer/service/buyer/runtime/config/chart"
	"github.com/xh3b4sd/wafer/service/buyer/runtime/config/surge"
)

type Config struct {
	Chart chart.Chart
	Surge surge.Surge
}
