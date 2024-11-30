package assert

import (
	"errors"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert/utils"
	"testing"
)

func Equal[V comparable](t *testing.T, expected, actual V) {
	t.Helper()

	if expected != actual {
		utils.Fatalf(t, actual, expected)
	}
}

func EqualSlice[V comparable](t *testing.T, expected, actual []V) {
	t.Helper()

	if len(expected) != len(actual) {
		utils.Fatalf(t, actual, expected)
	}

	for i := range expected {
		if expected[i] != actual[i] {
			utils.Fatalf(t, actual, expected)
		}
	}
}

func EqualSlice2D[V comparable](t *testing.T, expected, actual [][]V) {
	t.Helper()

	HasLen(t, len(expected), actual)
	for i, expectedLine := range expected {
		EqualSlice(t, expectedLine, actual[i])
	}
}

func EqualErr(t *testing.T, expected, actual error) {
	t.Helper()

	if expected == nil && actual != nil {
		utils.Fatalf(t, actual.Error(), expected)
	} else if actual == nil && expected != nil {
		utils.Fatalf(t, actual, expected.Error())
	} else if !errors.Is(actual, expected) {
		utils.Fatalf(t, actual.Error(), expected.Error())
	}
}
