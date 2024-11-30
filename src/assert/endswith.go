package assert

import (
	"fmt"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert/utils"
	"strings"
	"testing"
)

func EndsWith[V comparable](t *testing.T, expected string, actual V) {
	t.Helper()

	TypeOf[string](t, actual)

	if !strings.HasSuffix(any(actual).(string), expected) {
		utils.Fatalf(t, actual, fmt.Sprintf(".*%s", expected))
	}
}
