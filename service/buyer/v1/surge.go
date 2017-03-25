package v1

import (
	"math"
	"time"

	"github.com/xh3b4sd/wafer/service/informer"
)

type Surge struct {
	Angle      float64
	LeftBound  time.Time
	RightBound time.Time
}

// calculateSurge calculates the tangent of the imaginary rectengular triangle
// defined by the given coordinates. In case the surge is negative the angle is
// mutlipled by -1 to get its negative representation.
func calculateSurge(x1, y1, x2, y2 float64) float64 {
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

func calculateSurgeAverage(list []informer.Price) float64 {
	if len(list) < 2 {
		return 0
	}

	leftBound := list[0]
	rightBound := list[len(list)-1]

	surge := calculateSurge(float64(leftBound.Time.Unix()), leftBound.Buy, float64(rightBound.Time.Unix()), rightBound.Buy)

	return surge
}

func calculateSurgeDuration(list []informer.Price) time.Duration {
	if len(list) < 2 {
		return 0
	}

	leftBound := list[0]
	rightBound := list[len(list)-1]

	duration := rightBound.Time.Sub(leftBound.Time)

	return duration
}

func calculateSurgeTolerance(price, configTolerance float64) float64 {
	tolerance := price * configTolerance / 100
	return tolerance
}

func findLastSurge(prices []informer.Price, configTolerance float64) []informer.Price {
	if len(prices) < 2 {
		return []informer.Price{}
	}

	var n int
	var prevSurge informer.Price

	reversedSurges := reverse(prices)
	for i, s := range reversedSurges {
		tolerance := calculateSurgeTolerance(prevSurge.Buy, configTolerance)
		if i == 0 || prevSurge.Buy > (s.Buy-tolerance) {
			n = i
			prevSurge = s
			continue
		}

		break
	}
	lastSurges := reverse(reversedSurges[:n+1])

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
