package analyzer

import (
	"github.com/xh3b4sd/wafer/service/decider/config/analyzer/chart"
	"github.com/xh3b4sd/wafer/service/decider/config/analyzer/surge"
)

type Analyzer struct {
	Chart chart.Chart
	Surge surge.Surge
}
