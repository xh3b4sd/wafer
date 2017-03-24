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
	if r1.State.Trade.Revenue.Total != 0 {
		t.Fatal("expected", 0, "got", r1.State.Trade.Revenue.Total)
	}

	r1.State.Trade.Revenue.Total = 23.45

	r2 := tr.Runtime()
	if r2.State.Trade.Revenue.Total != 0 {
		t.Fatal("expected", 0, "got", r2.State.Trade.Revenue.Total)
	}
}
