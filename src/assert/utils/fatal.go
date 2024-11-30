package utils

import (
	"github.com/google/go-cmp/cmp"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/color"
	"testing"
)

func Fatalf(t *testing.T, got interface{}, want interface{}) {
	t.Helper()

	t.Fatalf("got %s%v%s, want %s%v%s\n%s", color.Red, got, color.Reset, color.Green, want, color.Reset, cmp.Diff(want, got))
}
