package config

import (
	"github.com/xh3b4sd/wafer/service/analyzer/iteration/config/analyzer"
	"github.com/xh3b4sd/wafer/service/analyzer/iteration/config/decider"
)

type Config struct {
	Analyzer analyzer.Analyzer
	Decider  decider.Decider
}
