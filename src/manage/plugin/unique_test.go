package plugin

import (
	"os"
	"path"
	"testing"
)

func TestDupGoFile(t *testing.T) {
	base := path.Join(os.Getenv("AOC_PUZZLES_DIR"), "2024-12-01/part1")

	cnf := &DupConfig{
		Base:        base,
        PackageName: "d01p01s01solution",
        TargetDir:   path.Join(base, "output", "1001-solution-tmp"),
    }

	files, err := os.ReadDir(base)
	if err != nil {
		t.Error(err)
	}

	for _, de := range files {
		if err := DupGoFile(cnf, de); err != nil {
			t.Error(err)
		}
	}
}

func TestCopyGoPackage(t *testing.T) {
	base := path.Join(os.Getenv("AOC_PUZZLES_DIR"), "2024-12-01/part1")

	cnf := &DupConfig{
		Base:        base,
        PackageName: "d01p01s01solution",
        TargetDir:   path.Join(base, "output", "1001-solution-tmp"),
    }

	if err := CopyGoPackage(cnf); err != nil {
		t.Error(err)
	}
}
