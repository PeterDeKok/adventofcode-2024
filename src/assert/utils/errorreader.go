package utils

import (
	"errors"
	"io"
)

type ErrorReader struct{}

var _ io.Reader = ErrorReader{}

func (e ErrorReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("forced error")
}
