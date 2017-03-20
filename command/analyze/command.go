// Package analyze implements the analyze command for the command line tool.
package analyze

import (
	"fmt"
	"os"
	"time"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	"github.com/spf13/cobra"

	"github.com/xh3b4sd/wafer/command/analyze/flag"
	"github.com/xh3b4sd/wafer/service/analyzer"
	"github.com/xh3b4sd/wafer/service/analyzer/iteration"
	"github.com/xh3b4sd/wafer/service/decider"
	"github.com/xh3b4sd/wafer/service/informer"
	"github.com/xh3b4sd/wafer/service/informer/csv"
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
	var err error

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
			return microerror.MaskAny(err)
		}
	}

	var newAnalyzer analyzer.Analyzer
	{
		config := iteration.DefaultConfig()
		config.Informer = newInformer
		config.Logger = c.logger
		newAnalyzer, err = iteration.New(config)
		if err != nil {
			return microerror.MaskAny(err)
		}
	}

	{
		initialConfig := decider.Config{}
		initialConfig.Analyzer.Chart.Window = 7 * 24 * time.Hour
		initialConfig.Analyzer.Surge.Min = 2.5
		initialConfig.Trader.Duration.Min = 10 * time.Minute
		initialConfig.Trader.Fee.Min = 4
		initialConfig.Trader.Revenue.Min = 2

		_, err := newAnalyzer.Iterate(initialConfig)
		if err != nil {
			return microerror.MaskAny(err)
		}

		fmt.Printf("%#v\n", newAnalyzer.Revenue())
	}

	return nil
}
