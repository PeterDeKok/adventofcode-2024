package assert

import (
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert/utils"
	"testing"
)

func HasLen[E any](t *testing.T, expected int, actual []E) {
	t.Helper()

	l := len(actual)

	if expected != l {
		utils.Fatalf(t, l, expected)
	}
}
