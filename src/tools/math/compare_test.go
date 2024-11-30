package math

import (
	"fmt"
	"math"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert/utils"
	"testing"
)

func TestNearlyEqual(t *testing.T) {
	values := []utils.Expected2[float64, float64, bool]{
		{ValueA: 0.0, ValueB: 0.0, Expected: true},
		{ValueA: 1.0, ValueB: 1.0, Expected: true},
		{ValueA: 999.0, ValueB: 999.0, Expected: true},
		{ValueA: 99999999.0, ValueB: 99999999.0, Expected: true},
		{ValueA: -1.0, ValueB: -1.0, Expected: true},
		{ValueA: -999.0, ValueB: -999.0, Expected: true},
		{ValueA: -99999999.0, ValueB: -99999999.0, Expected: true},
		{ValueA: -99999999.0 + math.SmallestNonzeroFloat64/2.0, ValueB: -99999999.0, Expected: true},
		{ValueA: 999.0 + math.SmallestNonzeroFloat64/2.0, ValueB: 999.0, Expected: true},
		{ValueA: 999.0 + 1e-13, ValueB: 999.0, Expected: false},
		{ValueA: -99999999.0, ValueB: -99999999.0 + math.SmallestNonzeroFloat64/2.0, Expected: true},
		{ValueA: 999.0, ValueB: 999.0 + math.SmallestNonzeroFloat64/2.0, Expected: true},
		{ValueA: 999.0, ValueB: 999.0 - 1e-13, Expected: false},
		{ValueA: -999.0, ValueB: 999.0, Expected: false},
		{ValueA: -0.0, ValueB: 0.0, Expected: true},
	}

	for _, exp := range values {
		t.Run(fmt.Sprintf("%f,%f", exp.ValueA, exp.ValueB), func(t *testing.T) {
			assert.Equal(t, exp.Expected, NearlyEqual(exp.ValueA, exp.ValueB))
		})
	}
}
