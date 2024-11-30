package math

import (
	"fmt"
	"math"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert/utils"
	"testing"
)

func TestAbsDiff(t *testing.T) {
	values := []utils.Expected2[int, int, int]{
		{ValueA: 0, ValueB: 0, Expected: 0},
		{ValueA: 0, ValueB: 1, Expected: 1},
		{ValueA: 0, ValueB: 10, Expected: 10},
		{ValueA: 0, ValueB: -1, Expected: 1},
		{ValueA: 0, ValueB: -10, Expected: 10},
		{ValueA: 0, ValueB: math.MinInt + 1, Expected: math.MaxInt},
		{ValueA: 0, ValueB: 0, Expected: 0},
		{ValueA: 1, ValueB: 0, Expected: 1},
		{ValueA: 10, ValueB: 0, Expected: 10},
		{ValueA: math.MaxInt, ValueB: 0, Expected: math.MaxInt},
		{ValueA: -1, ValueB: 0, Expected: 1},
		{ValueA: -10, ValueB: 0, Expected: 10},
		{ValueA: math.MinInt + 1, ValueB: 0, Expected: math.MaxInt},
		{ValueA: 1, ValueB: 1, Expected: 0},
		{ValueA: 10, ValueB: 10, Expected: 0},
		{ValueA: math.MaxInt, ValueB: math.MaxInt, Expected: 0},
		{ValueA: math.MinInt, ValueB: math.MinInt, Expected: 0},
		{ValueA: 1, ValueB: -1, Expected: 2},
		{ValueA: 10, ValueB: -10, Expected: 20},
		{ValueA: math.MaxInt / 2, ValueB: math.MinInt / 2, Expected: math.MaxInt},
		{ValueA: -1, ValueB: 1, Expected: 2},
		{ValueA: -10, ValueB: 10, Expected: 20},
		{ValueA: math.MinInt / 2, ValueB: math.MaxInt / 2, Expected: math.MaxInt},
	}

	for _, exp := range values {
		t.Run(fmt.Sprintf("%d,%d", exp.ValueA, exp.ValueB), func(t *testing.T) {
			assert.Equal(t, exp.Expected, AbsDiff(exp.ValueA, exp.ValueB))
		})
	}
}
