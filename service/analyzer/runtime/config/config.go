package config

import (
	"time"

	microerror "github.com/giantswarm/microkit/error"
	"github.com/spf13/cast"

	buyerconfig "github.com/xh3b4sd/wafer/service/buyer/runtime/config"
	permutationconfig "github.com/xh3b4sd/wafer/service/permutation/runtime/config"
	sellerconfig "github.com/xh3b4sd/wafer/service/seller/runtime/config"
	traderconfig "github.com/xh3b4sd/wafer/service/trader/runtime/config"
)

const (
	// Buyer.
	PermIDBuyerChartWindow      = "Buyer.Chart.Window"
	PermIDBuyerSurgeDurationMin = "Buyer.Surge.Duration.Min"
	PermIDBuyerSurgeMin         = "Buyer.Surge.Min"
	PermIDBuyerSurgeTolerance   = "Buyer.Surge.Tolerance"
	PermIDBuyerTradeCorridorMax = "Buyer.Trade.Corridor.Max"
	PermIDBuyerTradePauseMin    = "Buyer.Trade.Pause.Min"

	// Seller.
	PermIDSellerTradeDurationMin = "Seller.Trade.Duration.Min"

	// Trader.
	PermIDTraderTradeBudget = "Trader.Trade.Budget"
)

type Config struct {
	Buyer  buyerconfig.Config
	Seller sellerconfig.Config
	Trader traderconfig.Config
}

func (c *Config) GetPermConfigs() []permutationconfig.Config {
	var config permutationconfig.Config
	var configs []permutationconfig.Config

	//
	// Buyer.
	//

	//
	config = permutationconfig.Config{}
	config.ID = PermIDBuyerChartWindow
	config.Min = 24 * time.Hour
	config.Max = 24 * 30 * 12 * time.Hour
	config.Step = 24 * 30 * time.Hour
	configs = append(configs, config)

	//
	config = permutationconfig.Config{}
	config.ID = PermIDBuyerSurgeDurationMin
	config.Min = 10 * time.Second
	config.Max = 6 * time.Hour
	config.Step = 5 * time.Second
	configs = append(configs, config)

	//
	config = permutationconfig.Config{}
	config.ID = PermIDBuyerSurgeMin
	config.Min = 0.05
	config.Max = 5.0
	config.Step = 0.05
	configs = append(configs, config)

	//
	config = permutationconfig.Config{}
	config.ID = PermIDBuyerSurgeTolerance
	config.Min = 0.05
	config.Max = 5.0
	config.Step = 0.05
	configs = append(configs, config)

	//
	config = permutationconfig.Config{}
	config.ID = PermIDBuyerTradeCorridorMax
	config.Min = 75.0
	config.Max = 100.0
	config.Step = 0.5
	configs = append(configs, config)

	//
	config = permutationconfig.Config{}
	config.ID = PermIDBuyerTradePauseMin
	config.Min = 0
	config.Max = 24 * time.Hour
	config.Step = 1 * time.Hour
	configs = append(configs, config)

	//
	// Seller.
	//

	//
	config = permutationconfig.Config{}
	config.ID = PermIDSellerTradeDurationMin
	config.Min = 10 * time.Minute
	config.Max = 24 * 2 * time.Hour
	config.Step = 10 * time.Minute
	configs = append(configs, config)

	//
	// Trader.
	//

	//
	config = permutationconfig.Config{}
	config.ID = PermIDTraderTradeBudget
	config.Min = 100.0
	config.Max = 1000.0
	config.Step = 100.0
	configs = append(configs, config)

	return configs
}

func (c *Config) SetPermValue(permID string, permValue interface{}) error {
	switch permID {
	//
	// Buyer.
	//
	case PermIDBuyerChartWindow:
		d, err := cast.ToDurationE(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Buyer.Chart.Window = d
	case PermIDBuyerSurgeDurationMin:
		d, err := cast.ToDurationE(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Buyer.Surge.Duration.Min = d
	case PermIDBuyerSurgeMin:
		f, err := cast.ToFloat64E(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Buyer.Surge.Min = f
	case PermIDBuyerSurgeTolerance:
		f, err := cast.ToFloat64E(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Buyer.Surge.Tolerance = f
	case PermIDBuyerTradeCorridorMax:
		f, err := cast.ToFloat64E(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Buyer.Trade.Corridor.Max = f
	case PermIDBuyerTradePauseMin:
		d, err := cast.ToDurationE(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Buyer.Trade.Pause.Min = d
	//
	// Seller.
	//
	case PermIDSellerTradeDurationMin:
		d, err := cast.ToDurationE(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Seller.Trade.Duration.Min = d
	//
	// Trader.
	//
	case PermIDTraderTradeBudget:
		f, err := cast.ToFloat64E(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Trader.Trade.Budget = f
	default:
		return microerror.MaskAnyf(invalidExecutionError, "unknown permID '%s'", permID)
	}

	return nil
}

func (c *Config) Validate() error {
	var err error

	err = c.Buyer.Validate()
	if err != nil {
		return microerror.MaskAny(err)
	}
	err = c.Seller.Validate()
	if err != nil {
		return microerror.MaskAny(err)
	}
	err = c.Trader.Validate()
	if err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}
