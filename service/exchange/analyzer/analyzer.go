package analyzer

import (
	"bytes"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"
	chart "github.com/wcharczuk/go-chart"

	"github.com/xh3b4sd/wafer/service/exchange"
	"github.com/xh3b4sd/wafer/service/informer"
)

// Config is the configuration used to create a new exchange.
type Config struct {
	// Dependencies.
	Informer informer.Informer
	Logger   micrologger.Logger

	// Settings.
	BuyChan  chan informer.Price
	SellChan chan informer.Price
}

// DefaultConfig returns the default configuration used to create a new exchange
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Informer: nil,
		Logger:   nil,

		// Settings.
		BuyChan:  nil,
		SellChan: nil,
	}
}

// New creates a new configured exchange.
func New(config Config) (exchange.Exchange, error) {
	// Dependencies.
	if config.Informer == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Informer must not be empty")
	}
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

	newClient := &Exchange{
		// Dependencies.
		informer: config.Informer,
		logger:   config.Logger,

		// Internals.
		buyChan:  config.BuyChan,
		sellChan: config.SellChan,

		// Settings.
		buys:  []informer.Price{},
		sells: []informer.Price{},
	}

	return newClient, nil
}

// Exchange implements exchange.Exchange.
type Exchange struct {
	// Dependencies.
	informer informer.Informer
	logger   micrologger.Logger

	// Internals.
	buyChan  chan informer.Price
	sellChan chan informer.Price

	// Settings.
	buys  []informer.Price
	sells []informer.Price
}

func (e *Exchange) Buys() []informer.Price {
	return e.buys
}

func (e *Exchange) Execute() {
	for {
		select {
		case p := <-e.buyChan:
			e.buys = append(e.buys, p)
		case p := <-e.sellChan:
			e.sells = append(e.sells, p)
		}
	}
}

func (e *Exchange) Render() (*bytes.Buffer, error) {
	var xValues []float64
	var yValues []float64
	for p := range e.informer.Prices() {
		xValues = append(xValues, float64(p.Time.Unix()))
		yValues = append(yValues, p.Buy)
	}

	var gridLines []chart.GridLine
	for _, p := range e.buys {
		gl := chart.GridLine{
			Value: float64(p.Time.Unix()),
		}
		gridLines = append(gridLines, gl)
	}
	for _, p := range e.sells {
		gl := chart.GridLine{
			Value: float64(p.Time.Unix()),
		}
		gridLines = append(gridLines, gl)
	}

	graph := chart.Chart{
		Width: 1280,
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: chart.TimeHourValueFormatter,
			GridMajorStyle: chart.Style{
				Show:        true,
				StrokeColor: chart.ColorBlack,
				StrokeWidth: 0.2,
			},
			GridLines: gridLines,
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
		},
	}

	buf := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buf)
	if err != nil {
		return nil, microerror.MaskAny(err)
	}

	return buf, nil
}

func (e *Exchange) Sells() []informer.Price {
	return e.sells
}
