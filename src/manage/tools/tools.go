package tools

import (
	"errors"
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	ErrNotADir  = errors.New("not a directory")
	ErrNotAFile = errors.New("not a file")
)

func GetDirFromSrc(dir string) string {
	l := log.With("dir", dir)

	if _, f, _, ok := runtime.Caller(0); !ok {
		l.Fatal("failed to get puzzles directory", "reason", "root dir invalid")
	} else if absDir, err := filepath.Abs(filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(f))), dir)); err != nil {
		l.Fatal("failed to get puzzles directory", "err", err)
	} else if fi, err := os.Stat(absDir); err != nil {
		l.Fatal("failed to get puzzles directory", "err", err, "abs", absDir)
	} else if !fi.IsDir() {
		l.Fatal("failed to get puzzles directory", "reason", "not a directory")
	} else {
		return absDir
	}

	return ""
}

func FileExists(p string) (bool, error) {
	if fi, err := os.Stat(p); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else if fi.IsDir() {
		return false, ErrNotAFile
	}

	return true, nil
}

func DirExists(p string) (bool, error) {
	if fi, err := os.Stat(p); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else if !fi.IsDir() {
		return false, ErrNotADir
	}

	return true, nil
}
