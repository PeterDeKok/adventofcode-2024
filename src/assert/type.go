package assert

import (
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert/utils"
	"reflect"
	"testing"
)

func TypeOf[T any](t *testing.T, actual interface{}) {
	t.Helper()

	if _, ok := actual.(T); !ok {
		utils.Fatalf(t, reflect.TypeOf(actual), reflect.TypeFor[T]())
	}
}
