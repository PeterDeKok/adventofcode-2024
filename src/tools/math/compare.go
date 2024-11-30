package math

import "math"

func NearlyEqual(a, b float64) bool {
	if a == b {
		return true
	}

	d := math.Abs(a - b)

	if b == 0 {
		return d < math.SmallestNonzeroFloat64
	}

	return (d / math.Abs(b)) < math.SmallestNonzeroFloat64
}
