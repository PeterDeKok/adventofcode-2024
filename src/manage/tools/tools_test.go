package tools

import (
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"testing"
)

func TestDaysDir(t *testing.T) {
	assert.EndsWith(t, "/src/puzzles", GetDirFromSrc("puzzles"))
}
