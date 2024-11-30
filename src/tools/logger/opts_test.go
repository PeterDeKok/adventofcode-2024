package logger

import (
	"bytes"
	"io"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/assert"
	"testing"
)

func TestWithWriterOption(t *testing.T) {
	var buf io.Writer = bytes.NewBuffer([]byte{})
	il := &IterationLogger{}

	WithWriter(buf)(il)

	assert.Equal(t, buf, il.w)
}
