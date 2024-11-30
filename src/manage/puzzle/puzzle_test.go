package puzzle

import (
	"path"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/tools"
	"testing"
)

func TestValidatePuzzleDir(t *testing.T) {
	testDir := tools.GetDirFromSrc("../tests/fs/puzzles")

	tc := []struct {
		base    string
		dir     string
		wantOk  bool
		wantErr error
	}{
		{testDir, "not-exists", false, nil},
		{testDir, "dir-is-file", false, tools.ErrNotADir},
	}

	for _, tc := range tc {
		t.Run(tc.dir, func(t *testing.T) {
			ok, err := ValidatePuzzleDir(path.Join(tc.base, tc.dir))
			assert.Equal(t, tc.wantOk, ok)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
