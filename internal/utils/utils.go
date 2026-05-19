// Package utils provides shared helper utilities for gocron.
package utils

// RangeSlice returns a slice of integers from start (inclusive) to stop (exclusive).
func RangeSlice(start, stop int) []int {
	// Pre-allocate the exact size for efficiency
	out := make([]int, stop-start)
	for i := range out {
		out[i] = start + i
	}
	return out
}
