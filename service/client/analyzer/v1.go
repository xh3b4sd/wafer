package analyzer

import (
	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/xh3b4sd/wafer/service/client"
	"github.com/xh3b4sd/wafer/service/informer"
)

// Config is the configuration used to create a new client.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	BuyChan  chan informer.Price
	SellChan chan informer.Price
}

// DefaultConfig returns the default configuration used to create a new client
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		BuyChan:  nil,
		SellChan: nil,
	}
}

// New creates a new configured client.
func New(config Config) (client.Client, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Logger must not be empty")
	}

	// Settings.
	if config.BuyChan == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.BuyChan must not be empty")
	}
	if config.SellChan == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.SellChan must not be empty")
	}

	newClient := &Client{
		// Dependencies.
		logger: config.Logger,

		// Settings.
		buyChan:  config.BuyChan,
		sellChan: config.SellChan,
	}

	return newClient, nil
}

// Client implements client.Client.
type Client struct {
	// Dependencies.
	logger micrologger.Logger

	// Settings.
	buyChan  chan informer.Price
	sellChan chan informer.Price
}

func (c *Client) Buy(price informer.Price, volume float64) error {
	c.buyChan <- price
	return nil
}

func (c *Client) Close() error {
	close(c.buyChan)
	close(c.sellChan)
	return nil
}

func (c *Client) Sell(price informer.Price, volume float64) error {
	c.sellChan <- price
	return nil
}
