package v1

import (
	"math"
)

func calculateVolume(price, budget float64) float64 {
	f := budget / price
	v := roundDown(f, 2)

	return v
}

func roundDown(f float64, places int) float64 {
	p := math.Pow(10, float64(places))
	r := math.Floor(p * f)
	d := r / p

	return d
}
