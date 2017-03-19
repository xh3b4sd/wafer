package performance

import (
	"testing"
)

func Test_calculateSurge(t *testing.T) {
	testCases := []struct {
		X1       float64
		Y1       float64
		X2       float64
		Y2       float64
		Expected float64
	}{
		{
			X1:       0,
			Y1:       0,
			X2:       0,
			Y2:       0,
			Expected: 0,
		},
		{
			X1:       1,
			Y1:       1,
			X2:       1,
			Y2:       1,
			Expected: 0,
		},
		{
			X1:       3.432,
			Y1:       3.432,
			X2:       3.432,
			Y2:       3.432,
			Expected: 0,
		},
		{
			X1:       34.32,
			Y1:       34.32,
			X2:       34.32,
			Y2:       34.32,
			Expected: 0,
		},
		{
			X1:       1,
			Y1:       1,
			X2:       2,
			Y2:       2,
			Expected: 45,
		},
		{
			X1:       2,
			Y1:       2,
			X2:       1,
			Y2:       1,
			Expected: -45,
		},
		{
			X1:       1391212802,
			Y1:       797.4000000000,
			X2:       1391213041,
			Y2:       798.8900000000,
			Expected: 0.35719500199463483,
		},
	}

	for i, testCase := range testCases {
		expected := calculateSurge(testCase.X1, testCase.Y1, testCase.X2, testCase.Y2)
		if expected != testCase.Expected {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", expected)
		}
	}
}
