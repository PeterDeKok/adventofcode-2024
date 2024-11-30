package assert

import (
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert/utils"
	"testing"
)

func NoErr(t *testing.T, args ...any) {
	t.Helper()

	if len(args) == 0 {
		return
	}

	if err, ok := args[len(args)-1].(error); ok {
		utils.Fatalf(t, err, "no error")
	}
}
