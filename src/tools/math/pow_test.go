package math

import (
	"fmt"
	"math"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert/utils"
	"strconv"
	"testing"
)

func TestPowUint32(t *testing.T) {
	values := []utils.Expected2[uint32, uint32, uint32]{
		{ValueA: 0, ValueB: 0, Expected: 1},
		{ValueA: 12, ValueB: 0, Expected: 1},
		{ValueA: 12, ValueB: 1, Expected: 12},
		{ValueA: 12, ValueB: 2, Expected: 144},
		{ValueA: 12, ValueB: 5, Expected: 248832},
		{ValueA: 12, ValueB: 8, Expected: 429981696},
		{ValueA: 0, ValueB: 7, Expected: 0},
		{ValueA: 1, ValueB: 7, Expected: 1},
		{ValueA: 2, ValueB: 7, Expected: 128},
		{ValueA: 5, ValueB: 7, Expected: 78125},
		{ValueA: 8, ValueB: 7, Expected: 2097152},
		{ValueA: 2, ValueB: 31, Expected: 2147483648},
	}

	for _, exp := range values {
		t.Run(strconv.FormatUint(uint64(exp.ValueA), 10)+","+strconv.FormatUint(uint64(exp.ValueB), 10), func(t *testing.T) {
			assert.Equal(t, exp.Expected, PowUint32(exp.ValueA, exp.ValueB))
		})
	}
}

func TestIsPowerOfTwoUint32(t *testing.T) {
	values := []utils.Expected[uint32, bool]{
		{Value: 0, Expected: false},
		{Value: 3, Expected: false},
		{Value: 6, Expected: false},
		{Value: 7, Expected: false},
		{Value: 9, Expected: false},
		{Value: 10, Expected: false},
		{Value: math.MaxUint8, Expected: false},
		{Value: math.MaxUint16, Expected: false},
		{Value: math.MaxUint32, Expected: false},
	}

	for _, exp := range values {
		t.Run(strconv.FormatUint(uint64(exp.Value), 10), func(t *testing.T) {
			assert.Equal(t, exp.Expected, IsPowerOfTwoUint32(exp.Value))
		})
	}

	for i := uint32(0); i < 32; i++ {
		t.Run(fmt.Sprintf("2^%d", i), func(t *testing.T) {
			assert.Equal(t, true, IsPowerOfTwoUint32(PowUint32(2, i)))
		})
	}
}

func TestSqrtUint32(t *testing.T) {
	values := []utils.Expected[uint32, uint32]{
		{Value: 0, Expected: 0},
		{Value: 1, Expected: 1},
		{Value: 2, Expected: 1},
		{Value: 3, Expected: 1},
		{Value: 4, Expected: 2},
		{Value: 5, Expected: 2},
		{Value: math.MaxUint8, Expected: 15},
		{Value: math.MaxUint16, Expected: 255},
		{Value: math.MaxUint32, Expected: 65535},
	}

	for _, exp := range values {
		t.Run(strconv.FormatUint(uint64(exp.Value), 10), func(t *testing.T) {
			assert.Equal(t, exp.Expected, SqrtUint32(exp.Value))
		})
	}

	for i := uint32(1); i < 32; i++ {
		t.Run(fmt.Sprintf("sqrt(2^%d)", i), func(t *testing.T) {
			assert.Equal(t, i, SqrtUint32(PowUint32(i, 2)))
		})
	}

	for i := uint32(1); i < 32; i++ {
		t.Run(fmt.Sprintf("sqrt(2^%d)-1", i), func(t *testing.T) {
			assert.Equal(t, i-1, SqrtUint32(PowUint32(i, 2)-1))
		})
	}
}
