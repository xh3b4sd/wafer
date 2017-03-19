package performance

import (
	"fmt"
	"math"

	microerror "github.com/giantswarm/microkit/error"
	micrologger "github.com/giantswarm/microkit/logger"

	"github.com/xh3b4sd/wafer/service/decider"
	"github.com/xh3b4sd/wafer/service/informer"
)

// Config is the configuration used to create a new decider.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	DeciderConfig decider.Config
	Foo           string
}

// DefaultConfig returns the default configuration used to create a new decider
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		DeciderConfig: decider.Config{},
		Foo:           "",
	}
}

// New creates a new configured decider.
func New(config Config) (decider.Decider, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

	// Settings.
	if config.Foo == "" {
		return nil, microerror.MaskAnyf(invalidConfigError, "config.Foo must not be empty")
	}

	newAnalyzer := &Decider{
		// Internals.
		buyChan:  make(chan informer.Price, 10),
		buyEvent: true,
		buyPrice: informer.Price{},
		config:   config.DeciderConfig,
		sellChan: make(chan informer.Price, 10),
		window:   []informer.Price{},
	}

	return newAnalyzer, nil
}

// Decider implements decider.Decider.
type Decider struct {
	// Dependencies.
	logger micrologger.Logger

	// Internals.
	buyChan  chan informer.Price
	buyEvent bool
	buyPrice informer.Price
	config   decider.Config
	sellChan chan informer.Price
	window   []informer.Price
}

func (d *Decider) Buy() chan informer.Price {
	return d.buyChan
}

func (d *Decider) Config() decider.Config {
	return d.config
}

func (d *Decider) Sell() chan informer.Price {
	return d.sellChan
}

func (d *Decider) Watch(prices chan informer.Price) {
	for price := range prices {
		err := d.watch(price)
		if err != nil {
			d.logger.Log("error", fmt.Sprintf("%#v", err))
		}
	}
}

func (d *Decider) watch(price informer.Price) error {
	var err error
	d.window, err = calculateWindow(append(d.window, price), d.config.Analyzer.Chart.Window)
	if err != nil {
		return microerror.MaskAny(err)
	}

	views, err := calculateViews(d.window, d.config.Analyzer.Chart.View, d.config.Analyzer.Chart.Convolution)
	if IsNotEnoughData(err) {
		// In case there is not enough data yet, we cannot continue with the chart
		// analyzation. So we return here and wait for the next events and proceed
		// later, as soon as there is enough data for our algorithm.
		return nil
	} else if err != nil {
		return microerror.MaskAny(err)
	}

	if d.buyEvent {
		surges := viewsToSurges(views)
		surges = findLastSurge(surges)

		if surgeAverage(surges) < d.config.Analyzer.Surge.Min {
			return nil
		}

		d.buyChan <- price
		d.buyEvent = false
		d.buyPrice = price
	} else {
		revenue := calculateRevenue(d.buyPrice.Buy, price.Sell, d.config.Trader.Fee.Min)
		if revenue < d.config.Trader.Revenue.Min {
			return nil
		}

		d.sellChan <- price
		d.buyEvent = true
		d.buyPrice = informer.Price{}
	}

	return nil
}

// calculateRevenue takes the buy price and the current sell price to calculate
// the possible revenue with respect to some configured fee. The resulting
// floating point number is a percentage representation of the probable revenue
// based on the prise raise according to the original buy price.
func calculateRevenue(buyPrice, currentPrice, fee float64) float64 {
	total := currentPrice - buyPrice
	percentage := total * 100 / buyPrice
	withoutFee := percentage - fee

	if math.IsNaN(withoutFee) {
		return 0
	}

	return withoutFee
}

// TODO add configuration for tolerated surge burst
func findLastSurge(surges []Surge) []Surge {
	var n int
	var prevSurge Surge

	reversedSurges := reverse(surges)
	for i, s := range reversedSurges {
		if i == 0 || s.Angle < prevSurge.Angle {
			n = i
			prevSurge = s
			continue
		}

		break
	}
	reversedSurges = reversedSurges[:n+1]

	lastSurges := reverse(reversedSurges)

	if len(lastSurges) < 2 {
		return []Surge{}
	}

	return lastSurges
}

func reverse(list []Surge) []Surge {
	newList := make([]Surge, len(list))
	copy(newList, list)

	for i, j := 0, len(newList)-1; i < j; i, j = i+1, j-1 {
		newList[i], newList[j] = newList[j], newList[i]
	}

	return newList
}

func surgeAverage(list []Surge) float64 {
	var total float64

	for _, i := range list {
		total += i.Angle
	}

	average := total / float64(len(list))

	return average
}

func viewsToSurges(views []View) []Surge {
	var surges []Surge

	for _, v := range views {
		angle := calculateSurge(float64(v.LeftBound.Time.Unix()), v.LeftBound.Buy, float64(v.RightBound.Time.Unix()), v.RightBound.Buy)

		surge := Surge{
			Angle:      angle,
			LeftBound:  v.LeftBound.Time,
			RightBound: v.RightBound.Time,
		}

		surges = append(surges, surge)
	}

	return surges
}
