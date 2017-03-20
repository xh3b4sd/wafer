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
}

// DefaultConfig returns the default configuration used to create a new decider
// by best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		DeciderConfig: decider.Config{},
	}
}

// New creates a new configured decider.
func New(config Config) (decider.Decider, error) {
	// Dependencies.
	if config.Logger == nil {
		return nil, microerror.MaskAnyf(invalidConfigError, "logger must not be empty")
	}

	newAnalyzer := &Decider{
		// Dependencies.
		logger: config.Logger,

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

	d.window = append(d.window, price)
	d.window, err = calculateWindow(d.window, d.config.Analyzer.Chart.Window)
	if IsNotEnoughData(err) {
		// In case there is not enough data yet, we cannot continue with the chart
		// analyzation. So we return here and wait for the next events and proceed
		// later, as soon as there is enough data for our algorithm.
		return nil
	} else if err != nil {
		return microerror.MaskAny(err)
	}

	if d.buyEvent {
		prices := findLastSurge(d.window)
		surge := calculateSurgeAverage(prices)

		// TODO find out why surge is 2.5 and not 45
		if surge < d.config.Analyzer.Surge.Min {
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

		duration := price.Time.Sub(d.buyPrice.Time)
		if duration < d.config.Trader.Duration.Min {
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

func calculateSurgeAverage(list []informer.Price) float64 {
	if len(list) < 2 {
		return 0
	}

	leftBound := list[0]
	rightBound := list[len(list)-1]

	surge := calculateSurge(float64(leftBound.Time.Unix()), leftBound.Buy, float64(rightBound.Time.Unix()), rightBound.Buy)

	return surge
}

// TODO add configuration for tolerated surge burst
func findLastSurge(prices []informer.Price) []informer.Price {
	if len(prices) < 2 {
		return []informer.Price{}
	}

	var n int
	var prevSurge informer.Price

	reversedSurges := reverse(prices)
	for i, s := range reversedSurges {
		if i == 0 || s.Buy < prevSurge.Buy {
			n = i
			prevSurge = s
			continue
		}

		break
	}
	reversedSurges = reversedSurges[:n+1]

	lastSurges := reverse(reversedSurges)

	if len(lastSurges) < 2 {
		return []informer.Price{}
	}

	return lastSurges
}

func reverse(list []informer.Price) []informer.Price {
	newList := make([]informer.Price, len(list))
	copy(newList, list)

	for i, j := 0, len(newList)-1; i < j; i, j = i+1, j-1 {
		newList[i], newList[j] = newList[j], newList[i]
	}

	return newList
}
