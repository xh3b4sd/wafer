package v1

import (
	"math"
)

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
