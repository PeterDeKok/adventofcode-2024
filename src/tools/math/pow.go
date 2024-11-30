package math

import "math"

func PowUint32(n, m uint32) uint32 {
	if m == 0 {
		return 1
	}

	if m == 1 {
		return n
	}

	result := n
	for i := uint32(2); i <= m; i++ {
		result *= n
	}
	return result
}

func IsPowerOfTwoUint32(n uint32) bool {
	if n == 0 {
		return false
	}

	fn := float64(n)

	return NearlyEqual(math.Ceil(math.Log2(fn)), math.Floor(math.Log2(fn)))
}

func SqrtUint32(n uint32) uint32 {
	return uint32(math.Sqrt(float64(n)))
}
