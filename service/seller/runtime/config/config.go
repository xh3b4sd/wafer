package config

import (
	"time"

	microerror "github.com/giantswarm/microkit/error"
	"github.com/spf13/cast"

	permutationconfig "github.com/xh3b4sd/wafer/service/permutation/runtime/config"
	"github.com/xh3b4sd/wafer/service/seller/runtime/config/trade"
)

const (
	PermIDTradeDurationMin = "Trade.Duration.Min"
)

type Config struct {
	Trade trade.Trade
}

func (c *Config) GetPermConfigs() []permutationconfig.Config {
	var config permutationconfig.Config
	var configs []permutationconfig.Config

	//
	config = permutationconfig.Config{}
	config.ID = PermIDTradeDurationMin
	config.Min = 10 * time.Minute
	config.Max = 24 * 2 * time.Hour
	config.Step = 10 * time.Minute
	configs = append(configs, config)

	return configs
}

func (c *Config) SetPermValue(permID string, permValue interface{}) error {
	switch permID {
	case PermIDTradeDurationMin:
		d, err := cast.ToDurationE(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Trade.Duration.Min = d
	default:
		return microerror.MaskAnyf(invalidExecutionError, "unknown permID '%s'", permID)
	}

	return nil
}
