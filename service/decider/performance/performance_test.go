package performance

import (
	"reflect"
	"testing"
	"time"
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
		Surges   []Surge
		Expected []Surge
	}{
		{
			Surges: []Surge{
				{Angle: float64(10), LeftBound: time.Unix(1, 0), RightBound: time.Unix(2, 0)},
				{Angle: float64(20), LeftBound: time.Unix(3, 0), RightBound: time.Unix(4, 0)},
				{Angle: float64(30), LeftBound: time.Unix(5, 0), RightBound: time.Unix(6, 0)},
				{Angle: float64(40), LeftBound: time.Unix(7, 0), RightBound: time.Unix(8, 0)},
				{Angle: float64(10), LeftBound: time.Unix(9, 0), RightBound: time.Unix(10, 0)},
				{Angle: float64(20), LeftBound: time.Unix(11, 0), RightBound: time.Unix(12, 0)},
				{Angle: float64(30), LeftBound: time.Unix(13, 0), RightBound: time.Unix(14, 0)},
				{Angle: float64(40), LeftBound: time.Unix(15, 0), RightBound: time.Unix(16, 0)},
			},
			Expected: []Surge{
				{Angle: float64(10), LeftBound: time.Unix(9, 0), RightBound: time.Unix(10, 0)},
				{Angle: float64(20), LeftBound: time.Unix(11, 0), RightBound: time.Unix(12, 0)},
				{Angle: float64(30), LeftBound: time.Unix(13, 0), RightBound: time.Unix(14, 0)},
				{Angle: float64(40), LeftBound: time.Unix(15, 0), RightBound: time.Unix(16, 0)},
			},
		},
		{
			Surges: []Surge{
				{Angle: float64(10), LeftBound: time.Unix(1, 0), RightBound: time.Unix(2, 0)},
				{Angle: float64(20), LeftBound: time.Unix(3, 0), RightBound: time.Unix(4, 0)},
				{Angle: float64(30), LeftBound: time.Unix(5, 0), RightBound: time.Unix(6, 0)},
				{Angle: float64(40), LeftBound: time.Unix(7, 0), RightBound: time.Unix(8, 0)},
				{Angle: float64(10), LeftBound: time.Unix(9, 0), RightBound: time.Unix(10, 0)},
				{Angle: float64(25), LeftBound: time.Unix(11, 0), RightBound: time.Unix(12, 0)},
				{Angle: float64(50), LeftBound: time.Unix(13, 0), RightBound: time.Unix(14, 0)},
				{Angle: float64(90), LeftBound: time.Unix(15, 0), RightBound: time.Unix(16, 0)},
			},
			Expected: []Surge{
				{Angle: float64(10), LeftBound: time.Unix(9, 0), RightBound: time.Unix(10, 0)},
				{Angle: float64(25), LeftBound: time.Unix(11, 0), RightBound: time.Unix(12, 0)},
				{Angle: float64(50), LeftBound: time.Unix(13, 0), RightBound: time.Unix(14, 0)},
				{Angle: float64(90), LeftBound: time.Unix(15, 0), RightBound: time.Unix(16, 0)},
			},
		},
		{
			Surges: []Surge{
				{Angle: float64(10), LeftBound: time.Unix(1, 0), RightBound: time.Unix(2, 0)},
				{Angle: float64(10), LeftBound: time.Unix(3, 0), RightBound: time.Unix(4, 0)},
				{Angle: float64(10), LeftBound: time.Unix(5, 0), RightBound: time.Unix(6, 0)},
				{Angle: float64(10), LeftBound: time.Unix(7, 0), RightBound: time.Unix(8, 0)},
				{Angle: float64(10), LeftBound: time.Unix(9, 0), RightBound: time.Unix(10, 0)},
				{Angle: float64(10), LeftBound: time.Unix(11, 0), RightBound: time.Unix(12, 0)},
				{Angle: float64(50), LeftBound: time.Unix(13, 0), RightBound: time.Unix(14, 0)},
				{Angle: float64(90), LeftBound: time.Unix(15, 0), RightBound: time.Unix(16, 0)},
			},
			Expected: []Surge{
				{Angle: float64(10), LeftBound: time.Unix(11, 0), RightBound: time.Unix(12, 0)},
				{Angle: float64(50), LeftBound: time.Unix(13, 0), RightBound: time.Unix(14, 0)},
				{Angle: float64(90), LeftBound: time.Unix(15, 0), RightBound: time.Unix(16, 0)},
			},
		},
		{
			Surges: []Surge{
				{Angle: float64(90), LeftBound: time.Unix(1, 0), RightBound: time.Unix(2, 0)},
				{Angle: float64(80), LeftBound: time.Unix(3, 0), RightBound: time.Unix(4, 0)},
				{Angle: float64(70), LeftBound: time.Unix(5, 0), RightBound: time.Unix(6, 0)},
				{Angle: float64(60), LeftBound: time.Unix(7, 0), RightBound: time.Unix(8, 0)},
				{Angle: float64(50), LeftBound: time.Unix(9, 0), RightBound: time.Unix(10, 0)},
				{Angle: float64(40), LeftBound: time.Unix(11, 0), RightBound: time.Unix(12, 0)},
				{Angle: float64(10), LeftBound: time.Unix(13, 0), RightBound: time.Unix(14, 0)},
				{Angle: float64(10), LeftBound: time.Unix(15, 0), RightBound: time.Unix(16, 0)},
			},
			Expected: []Surge{},
		},
	}

	for i, testCase := range testCases {
		expected := findLastSurge(testCase.Surges)
		if !reflect.DeepEqual(expected, testCase.Expected) {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", expected)
		}
	}
}
