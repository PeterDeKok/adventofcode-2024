package assert

import "testing"

func ShouldPanic(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Got return, want panic")
		}
	}()

	fn()
}

func ShouldPanicWith(t *testing.T, fn func(), assert func(t, r any)) {
	t.Helper()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Got return, want panic")
		} else {
			assert(t, r)
		}
	}()

	fn()
}
