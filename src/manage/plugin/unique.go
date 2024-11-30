package plugin

import (
	"errors"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
)

var (
	ErrDupInputShort = errors.New("duplication input file contents too short")
	ErrDupInvalidInputPackageName = errors.New("duplication input file package name invalid")
)

type DupConfig struct {
	Base        string
	PackageName string
	TargetDir   string
}

func DupGoFile(cnf *DupConfig, de os.DirEntry) error {
	fn := de.Name()

	if de.IsDir() || !strings.HasSuffix(fn, ".go") || strings.HasSuffix(fn, "_test.go") {
		return nil
	}

	b, err := os.ReadFile(path.Join(cnf.Base, fn))
	if err != nil {
		return err
	}

	orig := []byte("package main\n")

	if len(b) < len(orig) {
		return ErrDupInputShort
	} else if !slices.Equal(b[:len(orig)], orig) {
		return ErrDupInvalidInputPackageName
	}

	pn := []byte(fmt.Sprintf("package main // %s\n", cnf.PackageName))
	trgt := make([]byte, len(pn)+len(b)-len(orig))

	copy(trgt[:len(pn)], pn[:])
	copy(trgt[len(pn):], b[len(orig):])

	return os.WriteFile(path.Join(cnf.TargetDir, fn), trgt, 0660)
}

func CopyGoPackage(cnf *DupConfig) error {
	dirEntries, err := os.ReadDir(cnf.Base)
	if err != nil {
		return err
	}

	if err := os.Mkdir(cnf.TargetDir, 0775); err != nil && !os.IsExist(err) {
		return err
	}

	for _, de := range dirEntries {
		if err := DupGoFile(cnf, de); err != nil {
			return err
		}
	}

	return nil
}
