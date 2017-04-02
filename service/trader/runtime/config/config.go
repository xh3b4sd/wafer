package config

import (
	microerror "github.com/giantswarm/microkit/error"
	"github.com/spf13/cast"

	permutationconfig "github.com/xh3b4sd/wafer/service/permutation/runtime/config"
	"github.com/xh3b4sd/wafer/service/trader/runtime/config/trade"
)

const (
	PermIDTradeBudget = "Trade.Budget"
)

type Config struct {
	Trade trade.Trade
}

func (c *Config) GetPermConfigs() []permutationconfig.Config {
	var config permutationconfig.Config
	var configs []permutationconfig.Config

	//
	config = permutationconfig.Config{}
	config.ID = PermIDTradeBudget
	config.Min = 100.0
	config.Max = 1000.0
	config.Step = 100.0
	configs = append(configs, config)

	return configs
}

func (c *Config) SetPermValue(permID string, permValue interface{}) error {
	switch permID {
	case PermIDTradeBudget:
		f, err := cast.ToFloat64E(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Trade.Budget = f
	default:
		return microerror.MaskAnyf(invalidExecutionError, "unknown permID '%s'", permID)
	}

	return nil
}

// TODO
func (c Config) Validate() error {
	return nil
}
