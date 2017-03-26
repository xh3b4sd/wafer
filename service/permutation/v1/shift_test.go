package v1

import (
	"reflect"
	"testing"
)

func Test_ShiftIndizes(t *testing.T) {
	testCases := []struct {
		Indizes  []int
		Max      []int
		Expected []int
	}{
		{
			Indizes:  []int{0, 0, 0},
			Max:      []int{3, 1, 2},
			Expected: []int{1, 0, 0},
		},
		{
			Indizes:  []int{1, 0, 0},
			Max:      []int{3, 1, 2},
			Expected: []int{2, 0, 0},
		},
		{
			Indizes:  []int{2, 0, 0},
			Max:      []int{3, 1, 2},
			Expected: []int{3, 0, 0},
		},
		{
			Indizes:  []int{3, 0, 0},
			Max:      []int{3, 1, 2},
			Expected: []int{0, 1, 0},
		},
		{
			Indizes:  []int{0, 1, 0},
			Max:      []int{3, 1, 2},
			Expected: []int{1, 1, 0},
		},
		{
			Indizes:  []int{1, 1, 0},
			Max:      []int{3, 1, 2},
			Expected: []int{2, 1, 0},
		},
		{
			Indizes:  []int{2, 1, 0},
			Max:      []int{3, 1, 2},
			Expected: []int{3, 1, 0},
		},
		{
			Indizes:  []int{3, 1, 0},
			Max:      []int{3, 1, 2},
			Expected: []int{0, 0, 1},
		},
		{
			Indizes:  []int{0, 0, 1},
			Max:      []int{3, 1, 2},
			Expected: []int{1, 0, 1},
		},
		{
			Indizes:  []int{1, 0, 1},
			Max:      []int{3, 1, 2},
			Expected: []int{2, 0, 1},
		},
		{
			Indizes:  []int{3, 1, 2},
			Max:      []int{3, 1, 2},
			Expected: []int{0, 0, 0},
		},
	}

	for i, testCase := range testCases {
		indizes, err := ShiftIndizes(testCase.Indizes, testCase.Max)
		if err != nil {
			t.Fatal("case", i+1, "expected", nil, "got", err)
		}

		if !reflect.DeepEqual(indizes, testCase.Expected) {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", indizes)
		}
	}
}

func Test_ShiftIndizes_Order(t *testing.T) {
	testCases := []struct {
		IndizesList [][]int
		Max         []int
	}{
		{
			IndizesList: [][]int{
				{0, 0, 0},
				{1, 0, 0},
				{2, 0, 0},
				{3, 0, 0},
				{0, 1, 0},
				{1, 1, 0},
				{2, 1, 0},
				{3, 1, 0},
				{0, 0, 1},
				{1, 0, 1},
				{2, 0, 1},
				{3, 0, 1},
				{0, 1, 1},
				{1, 1, 1},
				{2, 1, 1},
				{3, 1, 1},
				{0, 0, 2},
				{1, 0, 2},
				{2, 0, 2},
				{3, 0, 2},
				{0, 1, 2},
				{1, 1, 2},
				{2, 1, 2},
				{3, 1, 2},
			},
			Max: []int{3, 1, 2},
		},
	}

	for i, testCase := range testCases {
		for j, indizes := range testCase.IndizesList {
			if j == 0 {
				continue
			}

			from := testCase.IndizesList[j-1]
			to, err := ShiftIndizes(from, testCase.Max)
			if err != nil {
				t.Fatal("case", i+1, "iteration", j+1, "expected", nil, "got", err)
			}

			if !reflect.DeepEqual(to, indizes) {
				t.Fatal("case", i+1, "iteration", j+1, "expected", indizes, "got", to)
			}
		}
	}
}
