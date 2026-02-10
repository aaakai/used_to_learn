package utils

func MaxTwo[T int | float64 | int64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func MinTwo[T int | float64 | int64](a, b T) T {
	if a < b {
		return a
	}
	return b
}
