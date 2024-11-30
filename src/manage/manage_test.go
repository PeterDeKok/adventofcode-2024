package manage

import (
	"os"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/tools"
	"testing"
)

func TestCreate(t *testing.T) {
	m := Create(&AdventOfCodeConfig{
		PuzzlesDir: tools.GetDirFromSrc("../tests/fs/puzzles"),
		NrDays:     25,
		FirstDay:   "2024-12-01",
		TZ:         "Europe/Amsterdam",
	})

	assert.TypeOf[*Manager](t, m)
}

func TestManager_LoadLocal(t *testing.T) {
	m := Create(&AdventOfCodeConfig{
		PuzzlesDir: tools.GetDirFromSrc("../tests/fs/puzzles"),
		NrDays:     25,
		FirstDay:   "2024-12-01",
		TZ:         "Europe/Amsterdam",
	})

	m.LoadLocal()

	assert.NoErr(t, m.lookup["2024-12-01"].Error)
	assert.Equal(t, true, m.lookup["2024-12-01"].Part1.Validated)
	assert.NoErr(t, m.lookup["2024-12-01"].Part1.Error)
	assert.Equal(t, false, m.lookup["2024-12-01"].Part2.Validated)
	assert.EqualErr(t, os.ErrNotExist, m.lookup["2024-12-01"].Part2.Error)

	assert.NoErr(t, m.lookup["2024-12-02"].Error)
	assert.Equal(t, false, m.lookup["2024-12-02"].Part1.Validated)
	assert.EqualErr(t, os.ErrNotExist, m.lookup["2024-12-02"].Part1.Error)
	assert.Equal(t, false, m.lookup["2024-12-02"].Part2.Validated)
	assert.EqualErr(t, os.ErrNotExist, m.lookup["2024-12-02"].Part2.Error)
}
