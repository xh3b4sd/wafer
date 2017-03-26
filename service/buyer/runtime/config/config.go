package config

import (
	"time"

	microerror "github.com/giantswarm/microkit/error"
	"github.com/spf13/cast"

	"github.com/xh3b4sd/wafer/service/buyer/runtime/config/chart"
	"github.com/xh3b4sd/wafer/service/buyer/runtime/config/surge"
	"github.com/xh3b4sd/wafer/service/buyer/runtime/config/trade"
	permutationconfig "github.com/xh3b4sd/wafer/service/permutation/runtime/config"
)

const (
	PermIDChartWindow      = "Chart.Window"
	PermIDSurgeDurationMin = "Surge.Duration.Min"
	PermIDSurgeMin         = "Surge.Min"
	PermIDSurgeTolerance   = "Surge.Tolerance"
	PermIDTradePauseMin    = "Trade.Pause.Min"
)

type Config struct {
	Chart chart.Chart
	Surge surge.Surge
	Trade trade.Trade
}

func (c *Config) GetPermConfigs() []permutationconfig.Config {
	var config permutationconfig.Config
	var configs []permutationconfig.Config

	//
	config = permutationconfig.Config{}
	config.ID = PermIDChartWindow
	config.Min = 24 * time.Hour
	config.Max = 24 * 30 * 12 * time.Hour
	config.Step = 24 * 30 * time.Hour
	configs = append(configs, config)

	//
	config = permutationconfig.Config{}
	config.ID = PermIDSurgeDurationMin
	config.Min = 10 * time.Second
	config.Max = 6 * time.Hour
	config.Step = 5 * time.Second
	configs = append(configs, config)

	//
	config = permutationconfig.Config{}
	config.ID = PermIDSurgeMin
	config.Min = 0.05
	config.Max = 5.0
	config.Step = 0.05
	configs = append(configs, config)

	//
	config = permutationconfig.Config{}
	config.ID = PermIDSurgeTolerance
	config.Min = 0.05
	config.Max = 5.0
	config.Step = 0.05
	configs = append(configs, config)

	//
	config = permutationconfig.Config{}
	config.ID = PermIDTradePauseMin
	config.Min = 0
	config.Max = 24 * time.Hour
	config.Step = 1 * time.Hour
	configs = append(configs, config)

	return configs
}

func (c *Config) SetPermValue(permID string, permValue interface{}) error {
	switch permID {
	case PermIDChartWindow:
		d, err := cast.ToDurationE(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Chart.Window = d
	case PermIDSurgeDurationMin:
		d, err := cast.ToDurationE(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Surge.Duration.Min = d
	case PermIDSurgeMin:
		f, err := cast.ToFloat64E(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Surge.Min = f
	case PermIDSurgeTolerance:
		f, err := cast.ToFloat64E(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Surge.Tolerance = f
	case PermIDTradePauseMin:
		d, err := cast.ToDurationE(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Trade.Pause.Min = d
	default:
		return microerror.MaskAnyf(invalidExecutionError, "unknown permID '%s'", permID)
	}

	return nil
}
