// Package render implements the render command for the command line tool.
package render

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	"github.com/spf13/cobra"
	chart "github.com/wcharczuk/go-chart"

	"github.com/xh3b4sd/wafer/command/render/flag"
	"github.com/xh3b4sd/wafer/service/informer/csv"
)

var (
	f = &flag.Flag{}
)

// Config represents the configuration used to create a new render command.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger
}

// DefaultConfig provides a default configuration to create a new render command
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,
	}
}

// New creates a new configured render command.
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
		Use:   "render",
		Short: "Render charts as PNG files within the browser by running a simple server.",
		Long:  "Render charts as PNG files within the browser by running a simple server.",
		Run:   newCommand.Execute,
	}

	newCommand.cobraCommand.PersistentFlags().StringVar(&f.File, "file", "", "The absolute file path of a CSV file containing chart data.")
	newCommand.cobraCommand.PersistentFlags().BoolVar(&f.IgnoreHeader, "ignoreHeader", false, "Whether to ignore the first row of the CSV file.")
	newCommand.cobraCommand.PersistentFlags().IntVar(&f.Index.Buy, "index.buy", 0, "The index of the column within a CSV file representing buy prices.")
	newCommand.cobraCommand.PersistentFlags().IntVar(&f.Index.Sell, "index.sell", 0, "The index of the column within a CSV file representing sell prices.")
	newCommand.cobraCommand.PersistentFlags().IntVar(&f.Index.Time, "index.time", 0, "The index of the column within a CSV file representing price times.")
	newCommand.cobraCommand.PersistentFlags().IntVar(&f.Server.Port, "server.port", 8000, "The server port to listen on.")

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
	http.HandleFunc("/", drawChart)
	http.ListenAndServe(":"+strconv.Itoa(f.Server.Port), nil)

	return nil
}

func drawChart(res http.ResponseWriter, req *http.Request) {
	newConfig := csv.DefaultConfig()
	newConfig.BuyIndex = f.Index.Buy
	newConfig.File = f.File
	newConfig.IgnoreHeader = f.IgnoreHeader
	newConfig.SellIndex = f.Index.Sell
	newConfig.TimeIndex = f.Index.Time
	newInformer, err := csv.New(newConfig)
	if err != nil {
		panic(err)
	}

	var xValues []float64
	var yValues []float64
	for p := range newInformer.Prices() {
		xValues = append(xValues, float64(p.Time.Unix()))
		yValues = append(yValues, p.Buy)
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      "time in seconds",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Name:      "price in USD",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: xValues,
				YValues: yValues,
			},
			//			chart.AnnotationSeries{
			//				Annotations: []chart.Value2{
			//					{XValue: 1.0, YValue: 1.0, Label: "One"},
			//					{XValue: 2.0, YValue: 2.0, Label: "Two"},
			//					{XValue: 3.0, YValue: 3.0, Label: "Three"},
			//					{XValue: 4.0, YValue: 4.0, Label: "Four"},
			//					{XValue: 5.0, YValue: 5.0, Label: "Five"},
			//				},
			//			},
		},
	}

	res.Header().Set(http.CanonicalHeaderKey("Content-Type"), "image/png")
	graph.Render(chart.PNG, res)
}
