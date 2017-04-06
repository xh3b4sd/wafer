package v1

import (
	"testing"
)

// Test_Trader_Runtime_Copy makes sure the provided runtime information cannot
// be manipulated from the outside. That is, a copy of the runtime information
// is returned.
func Test_Trader_Runtime_Copy(t *testing.T) {
	tr := &Trader{}

	r1 := tr.Runtime()
	s1 := testSum(r1.State.Trade.Revenues)
	if s1 != 0 {
		t.Fatal("expected", 0, "got", s1)
	}

	r1.State.Trade.Revenues = []float64{23.45}

	r2 := tr.Runtime()
	s2 := testSum(r2.State.Trade.Revenues)
	if s2 != 0 {
		t.Fatal("expected", 0, "got", s2)
	}
}

func testSum(list []float64) float64 {
	var s float64

	for _, f := range list {
		s += f
	}

	return s
}
