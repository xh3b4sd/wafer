package performance

import (
	"math"
)

// Surge calculates the tangent of the imaginary rectengular triangle defined by
// the given coordinates. In case the surge is negative the angle is mutlipled
// by -1 to get its negative representation.
func Surge(x1, y1, x2, y2 float64) float64 {
	opposite := float64(y2 - y1)
	adjacent := float64(x2 - x1)
	atangent := math.Atan(opposite / adjacent)
	angle := atangent * 180 / math.Pi

	if math.IsNaN(angle) {
		return 0
	}

	if y2 < y1 {
		// The angle is negative. That means the price is declining.
		return angle * -1
	}

	return angle
}
