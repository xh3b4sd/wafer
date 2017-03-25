// Package analyze implements the analyze command for the command line tool.
package analyze

import (
	"fmt"
	"net/http"
	"os"
	"time"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/xh3b4sd/wafer/command/analyze/flag"
	"github.com/xh3b4sd/wafer/service/buyer"
	v1buyer "github.com/xh3b4sd/wafer/service/buyer/v1"
	"github.com/xh3b4sd/wafer/service/client"
	analyzerclient "github.com/xh3b4sd/wafer/service/client/analyzer"
	"github.com/xh3b4sd/wafer/service/exchange"
	analyzerexchange "github.com/xh3b4sd/wafer/service/exchange/analyzer"
	"github.com/xh3b4sd/wafer/service/informer"
	"github.com/xh3b4sd/wafer/service/informer/csv"
	"github.com/xh3b4sd/wafer/service/seller"
	v1seller "github.com/xh3b4sd/wafer/service/seller/v1"
	"github.com/xh3b4sd/wafer/service/trader"
	v1trader "github.com/xh3b4sd/wafer/service/trader/v1"
)

var (
	f = &flag.Flag{}
)

// Config represents the configuration used to create a new analyze command.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger
}

// DefaultConfig provides a default configuration to create a new analyze command
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,
	}
}

// New creates a new configured analyze command.
func New(config Config) (*Command, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

	newCommand := &Command{
		// Dependencies.
		logger: config.Logger,

		// Internals.
		cobraCommand: nil,
	}

	newCommand.cobraCommand = &cobra.Command{
		Use:   "analyze",
		Short: "Analyze charts to optimize the revenue produced by some algorithm.",
		Long:  "Analyze charts to optimize the revenue produced by some algorithm.",
		Run:   newCommand.Execute,
	}

	newCommand.cobraCommand.PersistentFlags().StringVar(&f.File, "file", "", "The absolute file path of a CSV file containing chart data.")
	newCommand.cobraCommand.PersistentFlags().BoolVar(&f.IgnoreHeader, "ignoreHeader", false, "Whether to ignore the first row of the CSV file.")
	newCommand.cobraCommand.PersistentFlags().IntVar(&f.Index.Buy, "index.buy", 0, "The index of the column within a CSV file representing buy prices.")
	newCommand.cobraCommand.PersistentFlags().IntVar(&f.Index.Sell, "index.sell", 0, "The index of the column within a CSV file representing sell prices.")
	newCommand.cobraCommand.PersistentFlags().IntVar(&f.Index.Time, "index.time", 0, "The index of the column within a CSV file representing price times.")

	return newCommand, nil
}

type Command struct {
	// Dependencies.
	logger micrologger.Logger

	// Internals.
	cobraCommand *cobra.Command
}

func (c *Command) CobraCommand() *cobra.Command {
	return c.cobraCommand
}

func (c *Command) Execute(cmd *cobra.Command, args []string) {
	err := f.Validate()
	if err != nil {
		c.logger.Log("error", fmt.Sprintf("%#v", microerror.MaskAny(err)))
		os.Exit(1)
	}

	err = c.execute()
	if err != nil {
		c.logger.Log("error", fmt.Sprintf("%#v", microerror.MaskAny(err)))
		os.Exit(1)
	}
}

func (c *Command) execute() error {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", c.handler).Methods("GET")

	http.Handle("/", rtr)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		return microerror.MaskAny(err)
	}

	return nil
}

func (c *Command) handler(res http.ResponseWriter, req *http.Request) {
	var err error

	buyChan := make(chan informer.Price, 100)
	sellChan := make(chan informer.Price, 100)

	var newInformer informer.Informer
	{
		config := csv.DefaultConfig()
		config.BuyIndex = f.Index.Buy
		config.File = f.File
		config.IgnoreHeader = f.IgnoreHeader
		config.SellIndex = f.Index.Sell
		config.TimeIndex = f.Index.Time
		newInformer, err = csv.New(config)
		if err != nil {
			panic(err)
		}
	}

	var newClient client.Client
	{
		config := analyzerclient.DefaultConfig()
		config.BuyChan = buyChan
		config.Logger = c.logger
		config.SellChan = sellChan
		newClient, err = analyzerclient.New(config)
		if err != nil {
			panic(err)
		}
		defer newClient.Close()
	}

	var newBuyer buyer.Buyer
	{
		config := v1buyer.DefaultConfig()
		config.Logger = c.logger
		config.Runtime.Chart.Window = 7 * 24 * time.Hour
		config.Runtime.Surge.Duration.Min = 20 * time.Minute
		config.Runtime.Surge.Min = 3.5
		config.Runtime.Surge.Tolerance = 0.6
		config.Runtime.Trade.Pause.Min = 5 * time.Hour
		newBuyer, err = v1buyer.New(config)
		if err != nil {
			panic(err)
		}
	}

	var newSeller seller.Seller
	{
		config := v1seller.DefaultConfig()
		config.Logger = c.logger
		config.Runtime.Chart.Window = 7 * 24 * time.Hour
		config.Runtime.Trade.Duration.Min = 4 * time.Minute
		config.Runtime.Trade.Fee.Min = 1.5
		config.Runtime.Trade.Revenue.Min = 4.5
		newSeller, err = v1seller.New(config)
		if err != nil {
			panic(err)
		}
	}

	var newTrader trader.Trader
	{
		config := v1trader.DefaultConfig()
		config.Buyer = newBuyer
		config.Client = newClient
		config.Informer = newInformer
		config.Logger = c.logger
		config.Seller = newSeller
		newTrader, err = v1trader.New(config)
		if err != nil {
			panic(err)
		}
	}

	var newExchange exchange.Exchange
	{
		config := analyzerexchange.DefaultConfig()
		config.BuyChan = buyChan
		config.Informer = newInformer
		config.Logger = c.logger
		config.SellChan = sellChan
		newExchange, err = analyzerexchange.New(config)
		if err != nil {
			panic(err)
		}
		defer newExchange.Close()
	}

	go func() {
		c.logger.Log("debug", "exchange started")
		newExchange.Execute()
		c.logger.Log("debug", "exchange stopped")
	}()

	c.logger.Log("debug", "trader started")
	err = newTrader.Execute()
	if err != nil {
		panic(err)
	}
	c.logger.Log("debug", "trader stopped")

	c.logger.Log("debug", "rendering started")
	buf, err := newExchange.Render()
	if err != nil {
		panic(err)
	}
	c.logger.Log("debug", "rendering finished")

	c.logger.Log("debug", "response started")
	res.Header().Set(http.CanonicalHeaderKey("Content-Type"), "image/png")
	res.Write(buf.Bytes())
	c.logger.Log("debug", "response finished")

	c.logger.Log("debug", fmt.Sprintf("trader revenue: %.2f", newTrader.Runtime().State.Trade.Revenue.Total))
}
