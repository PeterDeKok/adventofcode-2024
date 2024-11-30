package logger

import (
	"fmt"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert/utils"
	"testing"
)

func TestLevel_Fmt(t *testing.T) {
	values := []utils.Expected[Level, string]{
		{Value: LevelPanic, Expected: "\x1b[0;31mPAN\x1b[m"},
		{Value: LevelError, Expected: "\x1b[0;31mERR\x1b[m"},
		{Value: LevelWarn, Expected: "\x1b[0;33mWRN\x1b[m"},
		{Value: LevelInfo, Expected: "\x1b[0;34mINF\x1b[m"},
		{Value: LevelDebug, Expected: "\x1b[38;5;239mDBG\x1b[m"},
	}

	for _, exp := range values {
		t.Run(exp.Value.str, func(t *testing.T) {
			assert.Equal(t, exp.Expected, exp.Value.Fmt())
		})
	}
}

func TestLevel_Gt(t *testing.T) {
	levels := []Level{
		LevelPanic,
		LevelError,
		LevelWarn,
		LevelInfo,
		LevelDebug,
	}

	for i, lvl1 := range levels {
		for j, lvl2 := range levels {
			t.Run(fmt.Sprintf("%s vs %s", lvl1.str, lvl2.str), func(t *testing.T) {
				assert.Equal(t, i > j, lvl1.Gt(lvl2))
			})
		}
	}
}
