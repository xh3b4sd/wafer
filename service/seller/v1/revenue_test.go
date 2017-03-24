package v1

import (
	"testing"
)

func Test_calculateRevenue(t *testing.T) {
	testCases := []struct {
		BuyPrice  float64
		SellPrice float64
		Fee       float64
		Expected  float64
	}{
		// Test case 1 makes sure the zero value input causes a zero value result.
		{
			BuyPrice:  float64(0),
			SellPrice: float64(0),
			Fee:       float64(0),
			Expected:  float64(0),
		},
		// Test case 2 provides an example in which the price is raised by 10%. With
		// respect to a 5% fee the revenue is then 5%.
		{
			BuyPrice:  float64(100),
			SellPrice: float64(110),
			Fee:       float64(5),
			Expected:  float64(5),
		},
		// Test case 3 provides an example in which the price is declined by 10%. With
		// respect to a 5% fee the revenue is then -15%.
		{
			BuyPrice:  float64(100),
			SellPrice: float64(90),
			Fee:       float64(5),
			Expected:  float64(-15),
		},
	}

	for i, testCase := range testCases {
		expected := calculateRevenue(testCase.BuyPrice, testCase.SellPrice, testCase.Fee)
		if expected != testCase.Expected {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", expected)
		}
	}
}
