package assert

import (
	"fmt"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert/utils"
	"strings"
	"testing"
)

func StrCount[V comparable](t *testing.T, expected string, count int, actual V) {
	t.Helper()

	TypeOf[string](t, actual)

	if c := strings.Count(any(actual).(string), expected); c != count {
		utils.Fatalf(t, actual, fmt.Sprintf("count(%v) != %d", expected, count))
	}
}
