package performance

import (
	"reflect"
	"testing"
	"time"

	"github.com/xh3b4sd/wafer/service/informer"
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

func Test_findLastSurge(t *testing.T) {
	testCases := []struct {
		Prices   []informer.Price
		Expected []informer.Price
	}{
		{
			Prices:   []informer.Price{},
			Expected: []informer.Price{},
		},
		{
			Prices: []informer.Price{
				{Buy: 10.0, Sell: 10.0, Time: time.Unix(1, 0)},
				{Buy: 20.0, Sell: 20.0, Time: time.Unix(2, 0)},
				{Buy: 30.0, Sell: 30.0, Time: time.Unix(3, 0)},
				{Buy: 10.0, Sell: 10.0, Time: time.Unix(4, 0)},
				{Buy: 20.0, Sell: 20.0, Time: time.Unix(5, 0)},
				{Buy: 30.0, Sell: 30.0, Time: time.Unix(6, 0)},
			},
			Expected: []informer.Price{
				{Buy: 10.0, Sell: 10.0, Time: time.Unix(4, 0)},
				{Buy: 20.0, Sell: 20.0, Time: time.Unix(5, 0)},
				{Buy: 30.0, Sell: 30.0, Time: time.Unix(6, 0)},
			},
		},
		{
			Prices: []informer.Price{
				{Buy: 10.0, Sell: 10.0, Time: time.Unix(1, 0)},
				{Buy: 20.0, Sell: 20.0, Time: time.Unix(2, 0)},
				{Buy: 30.0, Sell: 30.0, Time: time.Unix(3, 0)},
				{Buy: 10.0, Sell: 10.0, Time: time.Unix(4, 0)},
				{Buy: 40.0, Sell: 40.0, Time: time.Unix(5, 0)},
				{Buy: 90.0, Sell: 90.0, Time: time.Unix(6, 0)},
			},
			Expected: []informer.Price{
				{Buy: 10.0, Sell: 10.0, Time: time.Unix(4, 0)},
				{Buy: 40.0, Sell: 40.0, Time: time.Unix(5, 0)},
				{Buy: 90.0, Sell: 90.0, Time: time.Unix(6, 0)},
			},
		},
		{
			Prices: []informer.Price{
				{Buy: 10.0, Sell: 10.0, Time: time.Unix(1, 0)},
				{Buy: 10.0, Sell: 20.0, Time: time.Unix(2, 0)},
				{Buy: 10.0, Sell: 30.0, Time: time.Unix(3, 0)},
				{Buy: 10.0, Sell: 10.0, Time: time.Unix(4, 0)},
				{Buy: 40.0, Sell: 40.0, Time: time.Unix(5, 0)},
				{Buy: 90.0, Sell: 90.0, Time: time.Unix(6, 0)},
			},
			Expected: []informer.Price{
				{Buy: 10.0, Sell: 10.0, Time: time.Unix(4, 0)},
				{Buy: 40.0, Sell: 40.0, Time: time.Unix(5, 0)},
				{Buy: 90.0, Sell: 90.0, Time: time.Unix(6, 0)},
			},
		},
		{
			Prices: []informer.Price{
				{Buy: 90.0, Sell: 90.0, Time: time.Unix(1, 0)},
				{Buy: 40.0, Sell: 40.0, Time: time.Unix(2, 0)},
				{Buy: 10.0, Sell: 10.0, Time: time.Unix(3, 0)},
				{Buy: 10.0, Sell: 20.0, Time: time.Unix(4, 0)},
				{Buy: 10.0, Sell: 30.0, Time: time.Unix(5, 0)},
				{Buy: 10.0, Sell: 10.0, Time: time.Unix(6, 0)},
			},
			Expected: []informer.Price{},
		},
	}

	for i, testCase := range testCases {
		expected := findLastSurge(testCase.Prices)
		if !reflect.DeepEqual(expected, testCase.Expected) {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", expected)
		}
	}
}
