package input

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"iter"
)

// CharReader returns an iterator which yields every character.
// It reads from the input until an error or [io.EOF] is encountered.
// A buffer is used to limit the impact of io operations.
// The bytes of the input are yielded one by one as a [rune].
//
// Multi-byte runes are not supported.
//
// The generator will panic if any unexpected errors are raised.
// This keeps the usage simple. When reading input results in an unexpected
// error, this will invalidate the entire input - and therefore run - regardless.
func CharReader(input io.Reader) iter.Seq2[int, rune] {
	return func(yield func(k int, v rune) bool) {
		rd := bufio.NewReader(input)

		for i := 0; ; i++ {
			if b, err := rd.ReadByte(); err != nil && errors.Is(err, io.EOF) {
				return
			} else if err != nil {
				panic(fmt.Errorf("failed to read bytes: %w", err))
			} else {
				yield(i, rune(b))
			}
		}
	}
}
