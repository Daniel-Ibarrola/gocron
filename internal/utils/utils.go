package utils

func RangeSlice(start, stop int) []int {
	// Pre-allocate the exact size for efficiency
	out := make([]int, stop-start)
	for i := range out {
		out[i] = start + i
	}
	return out
}
