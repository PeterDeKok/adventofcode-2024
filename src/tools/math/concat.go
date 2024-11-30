package math

// ToPow10 returns the next power of ten for a number.
// It can be used to 'shift' digits in base 10, which effectively
// allows the concatenation of the digits of two numbers.
//
// e.g.:
// [ToPow10](999) -> 1000
func ToPow10(n int) int {
	if n >= 1e18 {
		return 19
	}

	x := 10
	for x <= n {
		x *= 10
	}

	return x
}

// Concat concatenates the digits of two numbers
//
// Concatenation of digits can be seen as the first number,
// shifted - in base 10 - by the length of the second number,
// summed with the second number.
//
// e.g.:
// [Concat](1234, 56) -> 1234*100 + 56 -> 123456
func Concat(a, b int) int {
	return ToPow10(b)*a + b
}
