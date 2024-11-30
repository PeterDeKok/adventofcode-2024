package utils

type Expected[A comparable, V comparable] struct {
	Value    A
	Expected V
}

type Expected2[A comparable, B comparable, V comparable] struct {
	ValueA   A
	ValueB   B
	Expected V
}
