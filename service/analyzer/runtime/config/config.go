package config

import (
	"time"

	microerror "github.com/giantswarm/microkit/error"
	"github.com/spf13/cast"

	buyerconfig "github.com/xh3b4sd/wafer/service/buyer/runtime/config"
	permutationconfig "github.com/xh3b4sd/wafer/service/permutation/runtime/config"
	sellerconfig "github.com/xh3b4sd/wafer/service/seller/runtime/config"
)

const (
	// Buyer.
	PermIDBuyerTradeCorridorMax = "Buyer.Trade.Corridor.Max"
	PermIDBuyerTradePauseMin    = "Buyer.Trade.Pause.Min"

	// Seller.
	PermIDSellerTradeDurationMin = "Seller.Trade.Duration.Min"
	PermIDSellerTradeRevenueMin  = "Seller.Trade.Revenue.Min"
)

type Config struct {
	Buyer  buyerconfig.Config  `json:"buyer"`
	Seller sellerconfig.Config `json:"seller"`
}

func (c *Config) GetPermConfigs() []permutationconfig.Config {
	var config permutationconfig.Config
	var configs []permutationconfig.Config

	//
	// Buyer.
	//

	//
	config = permutationconfig.Config{}
	config.ID = PermIDBuyerTradeCorridorMax
	config.Min = 98.0
	config.Max = 100.0
	config.Step = 0.2
	configs = append(configs, config)

	//
	config = permutationconfig.Config{}
	config.ID = PermIDBuyerTradePauseMin
	config.Min = 2 * time.Hour
	config.Max = 4 * time.Hour
	config.Step = 15 * time.Minute
	configs = append(configs, config)

	//
	// Seller.
	//

	//
	config = permutationconfig.Config{}
	config.ID = PermIDSellerTradeDurationMin
	config.Min = 6 * time.Hour
	config.Max = 24 * time.Hour
	config.Step = 3 * time.Hour
	configs = append(configs, config)

	//
	config = permutationconfig.Config{}
	config.ID = PermIDSellerTradeRevenueMin
	config.Min = 2.0
	config.Max = 5.0
	config.Step = 0.2
	configs = append(configs, config)

	return configs
}

func (c *Config) SetPermValue(permID string, permValue interface{}) error {
	switch permID {
	//
	// Buyer.
	//
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
	case PermIDSellerTradeRevenueMin:
		f, err := cast.ToFloat64E(permValue)
		if err != nil {
			return microerror.MaskAny(err)
		}
		c.Seller.Trade.Revenue.Min = f
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

	return nil
}
