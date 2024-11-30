package assert

import (
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert/utils"
	"reflect"
	"testing"
)

type PT[V any] interface {
	*V
}

func NotNil(t *testing.T, actual interface{}) {
	t.Helper()

	if isNil(actual) {
		utils.Fatalf(t, actual, "!nil")
	}
}

// isNil checks if a specified object is nil or not, without Failing.
// https://github.com/stretchr/testify/blob/v1.10.0/assert/assertions.go#L674
func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)

	switch value.Kind() {
	case
		reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.UnsafePointer:

		return value.IsNil()
	default:
		return false
	}
}
