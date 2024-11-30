package input

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"iter"
	"strconv"
	"strings"
)

// LineReader returns an iterator which yields every line.
// It reads from the input until an error or [io.EOF] is encountered.
// A buffer is used to limit the impact of io operations.
//
// The standard buffer implementation optimises [bufio.Reader.ReadLine] operations
// around a max length by returning a partial line and an 'isPrefix' flag.
// This method does *not*! It will return the entire line, no matter the impact on performance!
//
// The generator will panic if any unexpected errors are raised.
// This keeps the usage simple. When reading input results in an unexpected
// error, this will invalidate the entire input - and therefore run - regardless.
func LineReader(input io.Reader) iter.Seq2[int, string] {
	return func(yield func(k int, v string) bool) {
		rd := bufio.NewReader(input)
		carry := strings.Builder{}

		for i := 0; ; {
			if line, isPrefix, err := rd.ReadLine(); err != nil && errors.Is(err, io.EOF) {
				if carry.Len() > 0 {
					yield(i, carry.String())
				}

				return
			} else if err != nil {
				panic(fmt.Errorf("failed to read lines: %w", err))
			} else if isPrefix {
				carry.Write(line)
				continue
			} else {
				carry.Write(line)

				if !yield(i, carry.String()) {
					return
				}

				carry.Reset()
				i++
			}
		}
	}
}

// LineIsStrSliceReader returns an iterator which yields a slice of strings for every line.
// The fields for a line are determined by [strings.Fields].
// i.e. This goal is reached by splitting on one or more values satisfying [unicode.IsSpace].
//
// The standard buffer implementation optimises [bufio.Reader.ReadLine] operations
// around a max length by returning a partial line and an 'isPrefix' flag.
// This method does *not*! It will return the entire line, no matter the impact on performance!
//
// The generator will panic if any unexpected errors are raised.
// This keeps the usage simple. When reading input results in an unexpected
// error, this will invalidate the entire input - and therefore run - regardless.
func LineIsStrSliceReader(input io.Reader) iter.Seq2[int, []string] {
	return func(yield func(k int, v []string) bool) {
		for i, line := range LineReader(input) {
			yield(i, strings.Fields(line))
		}
	}
}

// LineIsIntSliceReader returns an iterator which yields a slice of ints for every line.
// All fields in a line need to be convertable to integers or a panic is raised.
// This internally uses the [LineIsStrSliceReader] generator.
//
// The standard buffer implementation optimises [bufio.Reader.ReadLine] operations
// around a max length by returning a partial line and an 'isPrefix' flag.
// This method does *not*! It will return the entire line, no matter the impact on performance!
//
// The generator will panic if any unexpected errors are raised.
// This keeps the usage simple. When reading input results in an unexpected
// error, this will invalidate the entire input - and therefore run - regardless.
func LineIsIntSliceReader(input io.Reader) iter.Seq2[int, []int] {
	return func(yield func(k int, v []int) bool) {
		for i, tokens := range LineIsStrSliceReader(input) {
			numbers := make([]int, len(tokens))

			for j, token := range tokens {
				if n, err := strconv.Atoi(token); err != nil {
					panic(fmt.Errorf("failed to parse string %d to int on line %d: %w", j, i, err))
				} else {
					numbers[j] = n
				}
			}

			yield(i, numbers)
		}
	}
}
