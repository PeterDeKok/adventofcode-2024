package math

import (
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert/utils"
	"strconv"
	"testing"
)

func TestToPow10(t *testing.T) {
	values := []utils.Expected[int, int]{
		{Value: 0, Expected: 10},
		{Value: 1, Expected: 10},
		{Value: 2, Expected: 10},
		{Value: 4, Expected: 10},
		{Value: 7, Expected: 10},
		{Value: 9, Expected: 10},
		{Value: 10, Expected: 100},
		{Value: 11, Expected: 100},
		{Value: 20, Expected: 100},
		{Value: 99, Expected: 100},
		{Value: 100, Expected: 1000},
		{Value: 101, Expected: 1000},
		{Value: 106, Expected: 1000},
		{Value: 150, Expected: 1000},
		{Value: 999, Expected: 1000},
		{Value: 1000, Expected: 10000},
		{Value: 1001, Expected: 10000},
		{Value: 1006, Expected: 10000},
		{Value: 1500, Expected: 10000},
		{Value: 9999, Expected: 10000},
	}

	for _, exp := range values {
		t.Run(strconv.Itoa(exp.Value), func(t *testing.T) {
			assert.Equal(t, exp.Expected, ToPow10(exp.Value))
		})
	}
}

func TestConcat(t *testing.T) {
	values := []utils.Expected2[int, int, int]{
		{ValueA: 0, ValueB: 0, Expected: 0},
		{ValueA: 0, ValueB: 9, Expected: 9},
		{ValueA: 0, ValueB: 99, Expected: 99},
		{ValueA: 0, ValueB: 999, Expected: 999},
		{ValueA: 1, ValueB: 0, Expected: 10},
		{ValueA: 1, ValueB: 1, Expected: 11},
		{ValueA: 1, ValueB: 9, Expected: 19},
		{ValueA: 1, ValueB: 10, Expected: 110},
		{ValueA: 1, ValueB: 11, Expected: 111},
		{ValueA: 1, ValueB: 19, Expected: 119},
		{ValueA: 1, ValueB: 99, Expected: 199},
		{ValueA: 1, ValueB: 100, Expected: 1100},
		{ValueA: 1, ValueB: 109, Expected: 1109},
		{ValueA: 1, ValueB: 110, Expected: 1110},
		{ValueA: 1, ValueB: 111, Expected: 1111},
		{ValueA: 1, ValueB: 999, Expected: 1999},
		{ValueA: 9, ValueB: 0, Expected: 90},
		{ValueA: 9, ValueB: 1, Expected: 91},
		{ValueA: 9, ValueB: 9, Expected: 99},
		{ValueA: 9, ValueB: 10, Expected: 910},
		{ValueA: 9, ValueB: 11, Expected: 911},
		{ValueA: 9, ValueB: 19, Expected: 919},
		{ValueA: 9, ValueB: 99, Expected: 999},
		{ValueA: 9, ValueB: 111, Expected: 9111},
		{ValueA: 9, ValueB: 109, Expected: 9109},
		{ValueA: 9, ValueB: 999, Expected: 9999},
		{ValueA: 10, ValueB: 0, Expected: 100},
		{ValueA: 10, ValueB: 1, Expected: 101},
		{ValueA: 10, ValueB: 9, Expected: 109},
		{ValueA: 10, ValueB: 10, Expected: 1010},
		{ValueA: 10, ValueB: 99, Expected: 1099},
		{ValueA: 99, ValueB: 0, Expected: 990},
		{ValueA: 99, ValueB: 1, Expected: 991},
		{ValueA: 99, ValueB: 9, Expected: 999},
		{ValueA: 99, ValueB: 10, Expected: 9910},
		{ValueA: 99, ValueB: 99, Expected: 9999},
		{ValueA: 999, ValueB: 111, Expected: 999111},
		{ValueA: 4, ValueB: 294967295, Expected: 4294967295},
		{ValueA: 42949, ValueB: 67295, Expected: 4294967295},
		{ValueA: 429496729, ValueB: 5, Expected: 4294967295},
	}

	for _, exp := range values {
		t.Run(strconv.Itoa(exp.ValueA)+","+strconv.Itoa(exp.ValueB), func(t *testing.T) {
			assert.Equal(t, exp.Expected, Concat(exp.ValueA, exp.ValueB))
		})
	}
}
