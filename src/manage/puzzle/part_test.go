package puzzle

import (
	"fmt"
	"os"
	"path"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/tools"
	"testing"
)

func TestValidatePuzzlePartDir(t *testing.T) {
	testDir := tools.GetDirFromSrc("../tests/fs")

	tc := []struct {
		base    string
		dir     string
		wantOk  bool
		wantErr error
	}{
		{testDir, "not-exists", false, os.ErrNotExist},
		{testDir, "dir-is-file", false, tools.ErrNotADir},
		{testDir, "input-is-dir", false, tools.ErrNotAFile},
		{testDir, "no-input", false, os.ErrNotExist},
		{testDir, "no-problem", false, os.ErrNotExist},
		{testDir, "no-readme", false, os.ErrNotExist},
		{testDir, "no-sample-expected", false, os.ErrNotExist},
		{testDir, "no-sample-input", false, os.ErrNotExist},
		{testDir, "no-solution", false, os.ErrNotExist},
		{testDir, "no-solution-test", false, os.ErrNotExist},
		{testDir, "no-stats", false, os.ErrNotExist},
		{testDir, "problem-is-dir", false, tools.ErrNotAFile},
		{testDir, "readme-is-dir", false, tools.ErrNotAFile},
		{testDir, "sample-expected-is-dir", false, tools.ErrNotAFile},
		{testDir, "sample-input-is-dir", false, tools.ErrNotAFile},
		{testDir, "solution-is-dir", false, tools.ErrNotAFile},
		{testDir, "solution-test-is-dir", false, tools.ErrNotAFile},
		{testDir, "success", true, nil},
	}

	for _, tc := range tc {
		t.Run(tc.dir, func(t *testing.T) {
			r, err := ValidatePuzzlePartDir(nil, path.Join(tc.base, "puzzles", tc.dir))
			if build.DEBUG {
				fmt.Println(r)
			}
			assert.Equal(t, true, r.Done())
			assert.Equal(t, tc.wantOk, r.OK())
			assert.EqualErr(t, tc.wantErr, err)
			assert.EqualErr(t, tc.wantErr, r.Error())
		})
	}
}
