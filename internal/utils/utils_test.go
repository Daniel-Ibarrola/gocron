package utils

import (
	"reflect"
	"testing"
)

func TestRangeSlice(t *testing.T) {
	tests := []struct {
		name     string
		start    int
		stop     int
		expected []int
	}{
		{
			name:     "positive range",
			start:    1,
			stop:     5,
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "starts at zero",
			start:    0,
			stop:     3,
			expected: []int{0, 1, 2},
		},
		{
			name:     "single value range",
			start:    4,
			stop:     5,
			expected: []int{4},
		},
		{
			name:     "empty range when start equals stop",
			start:    3,
			stop:     3,
			expected: []int{},
		},
		{
			name:     "negative range",
			start:    -2,
			stop:     2,
			expected: []int{-2, -1, 0, 1},
		},
		{
			name:     "all negative range",
			start:    -5,
			stop:     -2,
			expected: []int{-5, -4, -3},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := RangeSlice(test.start, test.stop)

			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("RangeSlice(%d, %d) = %v, expected %v", test.start, test.stop, actual, test.expected)
			}
		})
	}
}
