package main

import (
	"testing"
)

type ShifterTest struct {
	N        int
	Values   []int
	Expected [][]int
}

func (s *ShifterTest) Match(result [][]int) bool {
	if len(result) != len(s.Expected) {
		return false
	}

	for i, row := range s.Expected {
		if len(result[i]) != len(row) {
			return false
		}

		for j, val := range row {
			if result[i][j] != val {
				return false
			}
		}
	}

	return true
}

func TestBitShifterShifting(t *testing.T) {
	tests := []ShifterTest{
		// First test
		ShifterTest{
			N:      2,
			Values: []int{1, 2, 3},
			Expected: [][]int{
				[]int{1, 2},
				[]int{1, 3},
				[]int{2, 3},
			},
		},
		// Second test
		ShifterTest{
			N:      2,
			Values: []int{1, 2, 3, 4},
			Expected: [][]int{
				[]int{1, 2},
				[]int{1, 3},
				[]int{1, 4},
				[]int{2, 3},
				[]int{2, 4},
				[]int{3, 4},
			},
		},
	}

	for _, test := range tests {
		bitShifter, _ := NewBitShifter(test.Values)
		combinations, _ := bitShifter.CombinationsN(test.N)
		if test.Match(combinations) == false {
			t.Errorf("BitShifter with Values: %+v expected %+v but got %+v", test.Values, test.Expected, combinations)
		}
	}
}
